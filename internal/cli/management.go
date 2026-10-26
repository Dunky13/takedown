package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/source"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type appRunner func(*cobra.Command, func(*app.App) error) error
type jobAccount interface {
	Jobs(context.Context) (source.JobSnapshot, error)
	RefreshJobs(context.Context) (source.JobSnapshot, error)
	RetryJob(context.Context, string) (string, error)
	CancelJob(context.Context, string) error
	ReconcileIntent(context.Context, string) error
	ResolveExpiredIntent(context.Context) error
	ArchiveJob(context.Context, string) error
	ResetAuthorization(context.Context) error
	CheckAccess(context.Context) (source.AccessInfo, error)
}

func confirmManagement(cmd *cobra.Command, confirmed bool, message string) error {
	if confirmed {
		return nil
	}
	input, ok := cmd.InOrStdin().(*os.File)
	if !ok || !term.IsTerminal(int(input.Fd())) {
		return errors.New("review the action and pass --confirm; no changes were made")
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "%s\nContinue? [y/N] ", message)
	answer, err := bufio.NewReader(input).ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(answer), "y") {
		return errors.New("action canceled")
	}
	return nil
}

func registerManagement(root *cobra.Command, run appRunner) {
	var online, probeStorage, doctorJSON bool
	doctor := &cobra.Command{Use: "doctor", Short: "Check setup; Google and storage probes are opt-in", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return run(cmd, func(a *app.App) error {
			r, err := a.Diagnose(cmd.Context(), app.DiagnosticOptions{Online: online, Storage: probeStorage, StoragePassword: os.Getenv("TAKEDOWN_STORAGE_PASSWORD")})
			if err != nil {
				return err
			}
			if doctorJSON {
				if err = json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
					return err
				}
			} else {
				for _, check := range r.Checks {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s: %s\n", strings.ToUpper(check.Status), check.Name, check.Message)
				}
			}
			if !r.Healthy() {
				return errors.New("one or more diagnostics failed; follow the displayed setup guidance")
			}
			return nil
		})
	}}
	doctor.Flags().BoolVar(&online, "online", false, "Use existing authorization for bounded Google discovery; never starts sign-in")
	doctor.Flags().BoolVar(&probeStorage, "storage", false, "Connect to destination and check free space; may create its directory, writes no archives")
	doctor.Flags().BoolVar(&doctorJSON, "json", false, "Print machine-readable diagnostic results")
	root.AddCommand(doctor)
	jobs := &cobra.Command{Use: "jobs", Short: "Inspect and explicitly manage recorded Data Portability jobs"}
	withJobs := func(cmd *cobra.Command, fn func(*app.App, jobAccount) error) error {
		return run(cmd, func(a *app.App) error {
			if a.Settings.EffectiveSource() != "portability" {
				return errors.New("job management requires --source portability and its matching OAuth profile")
			}
			m, ok := a.Account.(jobAccount)
			if !ok {
				return errors.New("this source does not support job management")
			}
			return fn(a, m)
		})
	}
	var refresh, jsonJobs bool
	list := &cobra.Command{Use: "list", Short: "Show pending, completed, failed, expired and archived jobs", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return withJobs(cmd, func(_ *app.App, m jobAccount) error {
			var snapshot source.JobSnapshot
			var err error
			if refresh {
				snapshot, err = m.RefreshJobs(cmd.Context())
			} else {
				snapshot, err = m.Jobs(cmd.Context())
			}
			if err != nil {
				return err
			}
			if jsonJobs {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(snapshot)
			}
			if snapshot.Intent != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "HELD %s: %s\n", snapshot.Intent.Kind, snapshot.Intent.Message)
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, "JOB\tSTATE\tACCESS\tRETRIES\tARCHIVED\tRESOURCES\tDETAILS")
			for _, j := range snapshot.Jobs {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%t\t%s\t%s\n", j.ID, j.State, j.AccessType, j.RetryCount, j.Archived, strings.Join(j.Resources, ","), j.Message)
			}
			return w.Flush()
		})
	}}
	list.Flags().BoolVar(&refresh, "refresh", false, "Refresh eligible jobs while respecting Google's polling cadence")
	list.Flags().BoolVar(&jsonJobs, "json", false, "Print credential-free job receipts")
	jobs.AddCommand(list)
	access := &cobra.Command{Use: "access", Short: "Ask Google's documented API which selected grants are one-time or time-based", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return withJobs(cmd, func(_ *app.App, m jobAccount) error {
			info, err := m.CheckAccess(cmd.Context())
			if err != nil {
				return err
			}
			return json.NewEncoder(cmd.OutOrStdout()).Encode(info)
		})
	}}
	jobs.AddCommand(access)
	type action struct {
		name, description, warning string
		args                       cobra.PositionalArgs
		apply                      func(context.Context, jobAccount, []string) (string, error)
	}
	actions := []action{
		{"retry JOB_ID", "Retry one failed job within Google's retry limit", "Starts a replacement Google export job. The request is never automatically replayed.", cobra.ExactArgs(1), func(ctx context.Context, m jobAccount, args []string) (string, error) {
			return m.RetryJob(ctx, args[0])
		}},
		{"cancel JOB_ID", "Cancel one in-progress time-based Google export", "Cancels this remote export job; local archives and other jobs remain.", cobra.ExactArgs(1), func(ctx context.Context, m jobAccount, args []string) (string, error) {
			return "Remote cancellation recorded.", m.CancelJob(ctx, args[0])
		}},
		{"reconcile JOB_ID", "Resolve a held action using an evidence-backed Google job ID", "Confirm that Google's receipt identifies the exact pending request, account, resources and creation time. The API can verify accessible job state, but cannot prove that association. Do not guess an ID; retain the hold if this evidence is unavailable.", cobra.ExactArgs(1), func(ctx context.Context, m jobAccount, args []string) (string, error) {
			return "Known job reconciled.", m.ReconcileIntent(ctx, args[0])
		}},
		{"resolve-expired", "Resolve an expired held intent only when retained evidence permits", "This does not bypass one-time grant or cadence guards. Unsafe or insufficient evidence keeps the hold.", cobra.NoArgs, func(ctx context.Context, m jobAccount, _ []string) (string, error) {
			return "Expired intent reconciled; safety history retained.", m.ResolveExpiredIntent(ctx)
		}},
		{"archive JOB_ID", "Archive eligible local history while retaining grant and cadence receipts", "Archives history only; it does not delete downloaded files, reset consent or cancel Google jobs.", cobra.ExactArgs(1), func(ctx context.Context, m jobAccount, args []string) (string, error) {
			return "Local job history archived; safety receipts retained.", m.ArchiveJob(ctx, args[0])
		}},
		{"reset-authorization", "Revoke Google consent after every known archive is downloaded", "Revokes Google authorization and permits removal of old remote archives. Only proceed after every known export and download is complete. Local files remain; sign in explicitly afterward.", cobra.NoArgs, func(ctx context.Context, m jobAccount, _ []string) (string, error) {
			return "Google authorization reset; explicitly authorize again before requesting another export.", m.ResetAuthorization(source.WithResetConfirmation(ctx))
		}},
	}
	for _, spec := range actions {
		spec := spec
		var confirmed bool
		command := &cobra.Command{Use: spec.name, Short: spec.description, Long: spec.description + "\n\n" + spec.warning, Args: spec.args, RunE: func(cmd *cobra.Command, args []string) error {
			if err := confirmManagement(cmd, confirmed, spec.warning); err != nil {
				return err
			}
			return withJobs(cmd, func(_ *app.App, m jobAccount) error {
				message, err := spec.apply(cmd.Context(), m, args)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), message)
				return nil
			})
		}}
		command.Flags().BoolVar(&confirmed, "confirm", false, "Confirm the described action; source safety checks still apply")
		jobs.AddCommand(command)
	}
	root.AddCommand(jobs)
}
