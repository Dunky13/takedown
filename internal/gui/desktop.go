//go:build !headless

// Package gui provides the native Fyne desktop interface.
package gui

import (
	"context"
	"errors"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"github.com/Dunky13/takedown/assets"
	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/config"
)

func Run(ctx context.Context, dataDir string, settings config.Settings) (retErr error) {
	a, err := app.New(dataDir, settings)
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, a.Close()) }()
	desktop := fyneapp.NewWithID("io.github.dunky13.takedown")
	desktop.SetIcon(fyne.NewStaticResource("icon.png", assets.IconPNG))
	metadata := desktop.Metadata()
	if metadata.ID == "" {
		metadata.ID = desktop.UniqueID()
	}
	if metadata.Migrations == nil {
		metadata.Migrations = make(map[string]bool)
	}
	metadata.Migrations["fyneDo"] = true
	fyneapp.SetMetadata(metadata)
	w := desktop.NewWindow("takedown")
	s := newScreen(ctx, desktop, w, a)
	defer s.close()
	w.Resize(fyne.NewSize(960, 760))
	w.SetContent(s.content)
	w.SetCloseIntercept(func() { s.cancel(); w.SetCloseIntercept(nil); w.Close() })
	go func() { <-s.ctx.Done(); fyne.Do(desktop.Quit) }()
	w.ShowAndRun()
	return nil
}
