package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Dunky13/takedown/internal/auth"
	"github.com/Dunky13/takedown/internal/source"
	"github.com/Dunky13/takedown/internal/storage"
)

type DiagnosticOptions struct {
	Online          bool
	Storage         bool
	StoragePassword string
}
type DiagnosticCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
type DiagnosticReport struct {
	Checks []DiagnosticCheck `json:"checks"`
}

func (r DiagnosticReport) Healthy() bool {
	for _, c := range r.Checks {
		if c.Status == "fail" {
			return false
		}
	}
	return true
}

// Diagnose checks local authorization without token refresh. Online discovery
// and destination connection are opt-in and never create exports. A storage
// check may create the configured destination directory, but no archive files.
func (a *App) Diagnose(ctx context.Context, opts DiagnosticOptions) (DiagnosticReport, error) {
	var report DiagnosticReport
	add := func(name, status, message string) {
		report.Checks = append(report.Checks, DiagnosticCheck{name, status, message})
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	if err := a.beginJob(); err != nil {
		return report, err
	}
	defer a.endJob()
	if err := a.Settings.Validate(); err != nil {
		add("preferences", "fail", err.Error())
		return report, nil
	}
	add("preferences", "pass", "Source and transfer settings are valid.")
	if a.Settings.EffectiveSource() == "drive" && a.Settings.DriveFolderID == "" {
		add("source selection", "fail", "Select the Takeout archive folder in Drive before discovery or watching.")
	}
	if a.Settings.EffectiveSource() == "portability" && len(a.Settings.PortabilityResources) == 0 {
		add("source selection", "fail", "Select only the Data Portability resource groups you intend to export before authorization.")
	}
	scopes, err := source.Scopes(a.Settings.EffectiveSource(), a.Settings.PortabilityResources)
	if err != nil {
		add("permissions", "fail", err.Error())
		return report, nil
	}
	m, err := auth.New(auth.Options{DataDir: a.DataDir, ClientFile: a.Settings.OAuthClientFile, Email: a.Settings.Email, Scopes: scopes})
	if err != nil {
		add("authorization", "fail", err.Error())
	} else {
		status, inspectErr := m.InspectAuthorization(ctx)
		_ = m.Close()
		if inspectErr != nil {
			add("authorization", "fail", inspectErr.Error())
		} else if !status.TokenStored {
			add("authorization", "warning", status.Message)
		} else if !status.RefreshAvailable && !status.AccessTokenExpiresAt.After(time.Now()) {
			add("authorization", "fail", "The saved access token expired and has no refresh grant; explicitly authorize again.")
		} else {
			add("authorization", "pass", status.Message)
		}
	}
	if opts.Online {
		checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		exports, checkErr := a.Account.List(checkCtx)
		cancel()
		if checkErr != nil {
			message := "Google API diagnostic failed; check authorization, selected resources, API configuration and connectivity."
			if _, safeSource := a.Account.(*source.Manager); safeSource {
				message = checkErr.Error()
			}
			add("Google API", "fail", message)
		} else {
			add("Google API", "pass", fmt.Sprintf("Discovery succeeded with %d ready archive entries. This does not verify the account hint or Google app approval.", len(exports)))
		}
	} else {
		add("Google API", "warning", "Not contacted. Enable online diagnostics to check authorization and source access.")
	}
	if opts.Storage {
		password, credentialErr := a.storagePassword(a.Settings.Destination, opts.StoragePassword)
		if credentialErr != nil {
			add("destination authorization", "fail", credentialErr.Error())
			return report, nil
		}
		checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		b, openErr := storage.Open(checkCtx, a.Settings.Destination, storage.Options{Password: password, PrivateKeyPath: a.Settings.PrivateKeyPath, KnownHostsPath: a.Settings.KnownHostsPath, AllowUnknownSpace: a.Settings.AllowUnknownSpace})
		if openErr != nil {
			add("destination", "fail", "Cannot open destination; check path, credentials, trusted host key and connectivity.")
		} else {
			free, spaceErr := b.FreeSpace(checkCtx)
			closeErr := b.Close()
			if spaceErr != nil {
				add("destination space", "fail", "Free space could not be verified; inspect server filesystem reporting and configured guard policy.")
			} else {
				add("destination space", "pass", fmt.Sprintf("Available: %d bytes; safety buffer: %d bytes. Archive writing and durable sync have not been tested.", free, a.Settings.SafetyBuffer))
			}
			if closeErr != nil {
				add("destination close", "fail", "Storage connection did not close cleanly.")
			}
		}
		cancel()
	} else {
		add("destination", "warning", "Not connected. A storage diagnostic may create the destination directory; it does not write archives.")
	}
	if a.Settings.EffectiveSource() == "portability" {
		add("Google policy", "warning", "Data Portability requires protected storage and applicable Google project eligibility/review. Encryption in transit and at rest is not verified by this application; establish it before downloading. No diagnostic result certifies compliance.")
	} else {
		add("Google policy", "warning", "Drive read-only access is broader than the configured folder. Applicable Google project review and public privacy-policy configuration are operator requirements, not locally verified.")
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	return report, nil
}
