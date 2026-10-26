package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type settingsReader func(*cobra.Command) (string, config.Settings, error)

func registerAutostart(root *cobra.Command, settings settingsReader) {
	service := &cobra.Command{Use: "service", Short: "Manage a watch in your OS login session; no system-wide service"}
	for _, action := range []string{"install", "remove", "status"} {
		action := action
		var dry, partByPart bool
		c := &cobra.Command{Use: action, Short: action + " this profile's user-login watch", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			dir, cfg, err := settings(cmd)
			if err != nil {
				return err
			}
			a, err := app.New(dir, cfg)
			if err != nil {
				return err
			}
			defer a.Close()
			result, err := a.ConfigureStartup(cmd.Context(), action, partByPart, dry)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
			return nil
		}}
		c.Flags().BoolVar(&dry, "dry-run", false, "Render the definition without writing files, reading tokens, or changing OS startup")
		c.Flags().BoolVar(&partByPart, "part-by-part", false, "Watch sequentially while space permits")
		service.AddCommand(c)
	}
	root.AddCommand(service)
	storage := &cobra.Command{Use: "storage", Short: "Manage destination passwords in your OS keyring"}
	store := &cobra.Command{Use: "auth", Short: "Save this SMB/SFTP destination password for unattended watching", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		dir, cfg, err := settings(cmd)
		if err != nil {
			return err
		}
		password := os.Getenv("TAKEDOWN_STORAGE_PASSWORD")
		if password == "" {
			in, ok := cmd.InOrStdin().(*os.File)
			if !ok || !term.IsTerminal(int(in.Fd())) {
				return errors.New("provide TAKEDOWN_STORAGE_PASSWORD privately, or run storage auth in a terminal for a hidden prompt")
			}
			fmt.Fprint(cmd.ErrOrStderr(), "Destination password: ")
			bytes, readErr := term.ReadPassword(int(in.Fd()))
			fmt.Fprintln(cmd.ErrOrStderr())
			if readErr != nil {
				return readErr
			}
			password = string(bytes)
			clear(bytes)
		}
		if err = config.SaveStoragePassword(dir, cfg.Destination, password); err != nil {
			return err
		}
		if err = config.Save(dir, cfg); err != nil {
			return errors.New("destination password saved, but preferences could not be saved; retain the same destination when retrying setup")
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Destination password saved in the OS keyring; it is not written into service files or preferences.")
		return nil
	}}
	forget := &cobra.Command{Use: "forget", Short: "Remove this destination's saved password", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		dir, cfg, err := settings(cmd)
		if err != nil {
			return err
		}
		if err = config.ForgetStoragePassword(dir, cfg.Destination); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Saved destination password removed; archives remain.")
		return nil
	}}
	storage.AddCommand(store, forget)
	root.AddCommand(storage)
}
