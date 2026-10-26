// Command releasectl validates commit conventions and plans native releases.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Dunky13/takedown/internal/releaseflow"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "releasectl:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("expected plan, check-title, check-commits, or validate-version")
	}
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	switch args[0] {
	case "plan":
		kind := flags.String("kind", "", "patch, minor, major, custom, or nightly")
		custom := flags.String("custom", "", "custom SemVer, without v")
		last := flags.String("last-nightly", "", "last successfully published nightly tag")
		output := flags.String("output", "", "new output directory")
		date := flags.String("date", "", "RFC3339 timestamp (default current UTC time)")
		if err := parse(flags, args[1:]); err != nil {
			return err
		}
		if *output == "" {
			return fmt.Errorf("--output is required")
		}
		now := time.Now().UTC()
		if *date != "" {
			var err error
			now, err = time.Parse(time.RFC3339, *date)
			if err != nil {
				return fmt.Errorf("invalid --date: %w", err)
			}
		}
		plan, notes, err := (releaseflow.Repository{Dir: "."}).Plan(releaseflow.Options{Kind: *kind, Custom: *custom, LastNightly: *last, Date: now})
		if err != nil {
			return err
		}
		return releaseflow.Write(*output, plan, notes)
	case "check-title":
		title := flags.String("title", "", "Conventional Commit or PR title")
		if err := parse(flags, args[1:]); err != nil {
			return err
		}
		return releaseflow.CheckTitle(*title)
	case "check-commits":
		base := flags.String("base", "", "exclusive base commit")
		head := flags.String("head", "HEAD", "inclusive head commit")
		if err := parse(flags, args[1:]); err != nil {
			return err
		}
		return (releaseflow.Repository{Dir: "."}).CheckCommits(*base, *head)
	case "validate-version":
		version := flags.String("version", "", "SemVer, without v")
		if err := parse(flags, args[1:]); err != nil {
			return err
		}
		_, err := releaseflow.ParseVersion(*version)
		return err
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func parse(flags *flag.FlagSet, args []string) error {
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	return nil
}
