package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/Dunky13/takedown/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), shutdownSignals()...)
	defer stop()
	if err := cli.Execute(ctx, os.Args[1:]); err != nil {
		if !errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, "takedown:", err)
		}
		os.Exit(1)
	}
}
