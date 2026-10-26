// Package cli exposes the same job flow as the native GUI through Cobra.
package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/gui"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/space"
	"github.com/Dunky13/takedown/internal/state"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var Version = "dev"

func Execute(ctx context.Context, args []string) error {
	c := NewCommand()
	c.SetArgs(args)
	return c.ExecuteContext(ctx)
}

func NewCommand() *cobra.Command {
	var dir, email, destination, key, hosts, provider, clientFile, driveFolder string
	var resources []string
	var interval time.Duration
	var workers int
	var buffer uint64
	var unknown bool
	root := &cobra.Command{Use: "takedown", Short: "Download and verify ready Google Takeout archives", Version: Version, SilenceUsage: true, SilenceErrors: true, Args: cobra.NoArgs}
	f := root.PersistentFlags()
	f.StringVar(&dir, "data-dir", "", "Private configuration and resumable state directory")
	f.StringVar(&email, "email", "", "Google account login hint and local OAuth profile slot")
	f.StringVar(&provider, "source", "drive", "Official export source: drive or portability")
	f.StringVar(&clientFile, "oauth-client", "", "Downloaded Google Desktop app OAuth client JSON file")
	f.StringVar(&driveFolder, "drive-folder", "", "Drive folder ID or link containing Takeout archives")
	f.StringSliceVar(&resources, "resources", nil, "Data Portability resource groups to authorize and export")
	f.DurationVar(&interval, "poll-interval", 5*time.Minute, "Watch interval (Drive >=1m, Data Portability >=5m)")
	f.StringVar(&destination, "destination", "", "Local, mounted, UNC, smb:// or sftp:// destination")
	f.StringVar(&key, "ssh-key", "", "SFTP private key file")
	f.StringVar(&hosts, "known-hosts", "", "SFTP known_hosts file")
	f.IntVar(&workers, "workers", 2, "Parallel archive downloads (1 to 16)")
	f.Uint64Var(&buffer, "safety-buffer", 2_000_000_000, "Minimum free bytes to retain")
	f.BoolVar(&unknown, "allow-unknown-space", false, "Allow SFTP without statvfs, with a visible warning")
	settings := func(cmd *cobra.Command) (string, config.Settings, error) {
		data := dir
		if data == "" {
			var err error
			data, err = config.DefaultDataDir()
			if err != nil {
				return "", config.Settings{}, err
			}
		}
		s, err := config.Load(data)
		if err != nil {
			return "", s, err
		}
		if cmd.Flags().Changed("email") {
			s.Email = email
		}
		if cmd.Flags().Changed("source") {
			s.Source = provider
		}
		if cmd.Flags().Changed("oauth-client") {
			s.OAuthClientFile = clientFile
		}
		if cmd.Flags().Changed("drive-folder") {
			s.DriveFolderID = driveFolder
		}
		if cmd.Flags().Changed("resources") {
			s.PortabilityResources = resources
		}
		if cmd.Flags().Changed("poll-interval") {
			if interval%time.Second != 0 || interval < 0 || interval > 24*time.Hour {
				return "", s, errors.New("poll interval must use whole seconds and be at most 24 hours")
			}
			s.PollIntervalSeconds = int(interval / time.Second)
		}
		if cmd.Flags().Changed("destination") {
			s.Destination = destination
		}
		if cmd.Flags().Changed("ssh-key") {
			s.PrivateKeyPath = key
		}
		if cmd.Flags().Changed("known-hosts") {
			s.KnownHostsPath = hosts
		}
		if cmd.Flags().Changed("workers") {
			s.Concurrency = workers
		}
		if cmd.Flags().Changed("safety-buffer") {
			s.SafetyBuffer = buffer
		}
		if cmd.Flags().Changed("allow-unknown-space") {
			s.AllowUnknownSpace = unknown
		}
		return data, s, s.Validate()
	}
	launch := func(cmd *cobra.Command, _ []string) error {
		data, s, err := settings(cmd)
		if err != nil {
			return err
		}
		return gui.Run(cmd.Context(), data, s)
	}
	root.RunE = launch
	root.AddCommand(&cobra.Command{Use: "gui", Short: "Open the native desktop application", Args: cobra.NoArgs, RunE: launch})
	withApp := func(cmd *cobra.Command, fn func(*app.App) error) (retErr error) {
		data, s, err := settings(cmd)
		if err != nil {
			return err
		}
		a, err := app.New(data, s)
		if err != nil {
			return err
		}
		defer func() { retErr = errors.Join(retErr, a.Close()) }()
		stop := subscribe(a.Bus, cmd.ErrOrStderr())
		defer stop()
		return fn(a)
	}
	var forget bool
	authCmd := &cobra.Command{Use: "auth", Short: "Authorize official APIs in your system browser", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return withApp(cmd, func(a *app.App) error {
			if forget {
				if err := a.Account.Forget(cmd.Context()); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Local authorization removed for this profile. Revoke Google's remote grant separately at https://myaccount.google.com/permissions.")
				return nil
			}
			if a.Settings.EffectiveSource() == "drive" {
				fmt.Fprintln(cmd.ErrOrStderr(), "Authorization requests read-only access to view and download Drive files. takedown scans only your configured folder.")
			} else {
				fmt.Fprintln(cmd.ErrOrStderr(), "Authorization requests export access only for your selected Data Portability resources.")
				fmt.Fprintln(cmd.ErrOrStderr(), "Data Portability archives require encryption in transit and at rest. takedown does not encrypt destination files or verify your storage's encryption. Establish protected storage before downloading: SFTP encrypts transit only; SMB encryption depends on the server and share configuration.")
			}
			fmt.Fprintln(cmd.ErrOrStderr(), "OAuth tokens are stored in your OS keyring. Archives go to your chosen destination; no data or telemetry is sent to the developer. Disconnect with auth --forget; archive deletion remains under your control. Complete Google's consent in your browser.")
			if err := a.Account.Login(cmd.Context()); err != nil {
				return err
			}
			if err := config.Save(a.DataDir, a.Settings); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "OAuth authorization saved in the OS keyring. Passwords and MFA stay in your browser.")
			return nil
		})
	}}
	authCmd.Flags().BoolVar(&forget, "forget", false, "Disconnect and remove this profile's saved OAuth authorization")
	root.AddCommand(authCmd)
	var asJSON bool
	listCmd := &cobra.Command{Use: "list", Short: "Discover ready Takeout exports and exact part sizes", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return withApp(cmd, func(a *app.App) error {
			exports, err := a.Account.List(cmd.Context())
			if err != nil {
				return err
			}
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(exports)
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, "EXPORT\tCREATED\tPART\tBYTES\tARCHIVE")
			for _, e := range exports {
				for _, p := range e.Parts {
					fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n", e.ID, e.CreatedAt.Format("2006-01-02 15:04"), p.Index, p.Size, p.Filename)
				}
			}
			if err = w.Flush(); err != nil {
				return err
			}
			if len(exports) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No ready exports found.")
			}
			return nil
		})
	}}
	listCmd.Flags().BoolVar(&asJSON, "json", false, "Print machine-readable export metadata")
	root.AddCommand(listCmd)
	var exportWatch bool
	exportCmd := &cobra.Command{Use: "export", Short: "Open Google Takeout setup, or start an authorized Data Portability export", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return withApp(cmd, func(a *app.App) error {
			result, err := a.Account.CreateExport(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
			if err = config.Save(a.DataDir, a.Settings); err != nil {
				return err
			}
			if exportWatch {
				return a.Watch(cmd.Context(), app.WatchOptions{DownloadOptions: app.DownloadOptions{StoragePassword: os.Getenv("TAKEDOWN_STORAGE_PASSWORD")}})
			}
			return nil
		})
	}}
	exportCmd.Flags().BoolVar(&exportWatch, "watch", false, "After export setup, wait for and download ready archives until canceled")
	root.AddCommand(exportCmd)
	var selected []string
	var sequential bool
	downloadCmd := &cobra.Command{Use: "download", Short: "Download ready archives, resuming existing partial files", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return withApp(cmd, func(a *app.App) error {
			exports, err := a.Account.List(cmd.Context())
			if err != nil {
				return err
			}
			parts, err := SelectParts(exports, selected)
			if err != nil {
				return err
			}
			opts := app.DownloadOptions{Destination: a.Settings.Destination, PartByPart: sequential, StoragePassword: os.Getenv("TAKEDOWN_STORAGE_PASSWORD")}
			err = a.Download(cmd.Context(), parts, opts)
			if errors.Is(err, space.ErrInsufficientSpace) && !opts.PartByPart {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				if !term.IsTerminal(int(os.Stdin.Fd())) {
					return errors.New("batch does not fit: free space, change --destination, or explicitly use --part-by-part")
				}
				fmt.Fprint(cmd.ErrOrStderr(), "Download one part at a time while space permits? [y/N] ")
				answer, readErr := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
				if readErr != nil {
					return readErr
				}
				if strings.EqualFold(strings.TrimSpace(answer), "y") {
					opts.PartByPart = true
					return a.Download(cmd.Context(), parts, opts)
				}
			}
			return err
		})
	}}
	downloadCmd.Flags().StringSliceVar(&selected, "export", nil, "Export IDs to download (default: all ready exports)")
	downloadCmd.Flags().BoolVar(&sequential, "part-by-part", false, "Download sequentially, pausing when the next part cannot fit")
	root.AddCommand(downloadCmd)
	watchCmd := &cobra.Command{Use: "watch", Short: "Continuously download new archives through official APIs", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return withApp(cmd, func(a *app.App) error {
			return a.Watch(cmd.Context(), app.WatchOptions{DownloadOptions: app.DownloadOptions{Destination: a.Settings.Destination, PartByPart: sequential, StoragePassword: os.Getenv("TAKEDOWN_STORAGE_PASSWORD")}, ExportIDs: selected})
		})
	}}
	watchCmd.Flags().StringSliceVar(&selected, "export", nil, "Watch only these export IDs (default: all available archives)")
	watchCmd.Flags().BoolVar(&sequential, "part-by-part", false, "Download sequentially, pausing when the next part cannot fit")
	root.AddCommand(watchCmd)
	var statusJSON bool
	statusCmd := &cobra.Command{Use: "status", Short: "Read durable download progress without opening a browser", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		data, _, err := settings(cmd)
		if err != nil {
			return err
		}
		records, err := state.ReadSnapshot(filepath.Join(data, "takedown-state.json"))
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(cmd.OutOrStdout(), "No downloads recorded.")
			return nil
		}
		if err != nil {
			return err
		}
		if statusJSON {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(records)
		}
		sort.Slice(records, func(i, j int) bool { return records[i].Part.Key() < records[j].Part.Key() })
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "PART\tBYTES\tTOTAL\tCOMPLETE\tDESTINATION")
		for _, r := range records {
			fmt.Fprintf(w, "%s\t%d\t%d\t%t\t%s\n", r.Part.Key(), r.Offset, r.Part.Size, r.Complete, r.Destination)
		}
		return w.Flush()
	}}
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "Print machine-readable durable progress")
	root.AddCommand(statusCmd)
	registerManagement(root, withApp)
	registerAutostart(root, settings)
	return root
}

func SelectParts(exports []model.Export, ids []string) ([]model.Part, error) {
	wanted := make(map[string]bool)
	for _, id := range ids {
		wanted[id] = false
	}
	var parts []model.Part
	for _, e := range exports {
		if len(wanted) > 0 {
			if _, ok := wanted[e.ID]; !ok {
				continue
			}
			wanted[e.ID] = true
		}
		parts = append(parts, e.Parts...)
	}
	for id, found := range wanted {
		if !found {
			return nil, fmt.Errorf("ready export not found: %s", id)
		}
	}
	if len(parts) == 0 {
		return nil, errors.New("no ready archive parts found")
	}
	return parts, nil
}
func subscribe(bus *events.Bus, w io.Writer) func() {
	ch, cancel := bus.Subscribe(256)
	done := make(chan struct{})
	var once sync.Once
	go func() {
		defer close(done)
		bars := make(map[string]*progressbar.ProgressBar)
		terminal := false
		if f, ok := w.(*os.File); ok {
			terminal = term.IsTerminal(int(f.Fd()))
		}
		for e := range ch {
			if e.Kind == events.Progress {
				if terminal {
					bar := bars[e.PartKey]
					if bar == nil {
						label := e.Filename
						if label == "" {
							label = "Archive"
						}
						bar = progressbar.NewOptions64(e.Total, progressbar.OptionSetWriter(w), progressbar.OptionSetDescription(label), progressbar.OptionShowBytes(true), progressbar.OptionThrottle(250_000_000))
						bars[e.PartKey] = bar
					}
					_ = bar.Set64(e.Bytes)
				}
				continue
			}
			if e.Kind == events.PartCompleted {
				if b := bars[e.PartKey]; b != nil {
					_ = b.Finish()
					delete(bars, e.PartKey)
				}
			}
			fmt.Fprintf(w, "%s: %s\n", e.Kind, e.Message)
		}
	}()
	return func() { once.Do(func() { cancel(); <-done }) }
}
