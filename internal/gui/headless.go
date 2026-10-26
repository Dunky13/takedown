//go:build headless

package gui

import (
	"context"
	"errors"
	"github.com/Dunky13/takedown/internal/config"
)

func Run(context.Context, string, config.Settings) error {
	return errors.New("this optional headless build has no GUI; install the default desktop binary or run a CLI subcommand")
}
