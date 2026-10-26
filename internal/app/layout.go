package app

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/state"
)

func layoutParts(parts []model.Part, destination string, records []state.Record) []model.Part {
	byKey := make(map[string]state.Record, len(records))
	for _, record := range records {
		byKey[record.Part.Key()] = record
	}
	out := append([]model.Part(nil), parts...)
	for i := range out {
		p := &out[i]
		if old, ok := byKey[p.Key()]; ok && old.Destination == destination {
			// Retain an established output path, including pre-namespace flat
			// checkpoints. The engine separately verifies all immutable metadata.
			p.RelativePath = old.Part.RelativePath
		} else {
			p.RelativePath = model.NamespacePath(p.ExportID, p.Filename)
		}
	}
	return out
}

func (a *App) previewLayout(parts []model.Part, destination string) ([]model.Part, error) {
	id, err := DestinationIdentity(destination)
	if err != nil {
		return nil, err
	}
	records, err := state.ReadSnapshot(filepath.Join(a.DataDir, "takedown-state.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return layoutParts(parts, id, records), nil
}
