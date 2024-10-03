package txtarfmt

import (
	"bytes"
	"encoding/json"
	"go/format"
	"path/filepath"

	"golang.org/x/tools/txtar"
)

type Configuration struct {
	SkipGo, SkipJSON bool
}

func Archive(archive *txtar.Archive, config Configuration) error {
	for i, file := range archive.Files {
		fmtFile, err := File(file, config)
		if err != nil {
			return err
		}
		archive.Files[i] = fmtFile
	}
	return nil
}

func File(file txtar.File, config Configuration) (txtar.File, error) {
	if !config.SkipGo && filepath.Ext(file.Name) == ".go" {
		out, err := format.Source(file.Data)
		if err != nil {
			return file, err
		}
		file.Data = out
	}
	if !config.SkipJSON && filepath.Ext(file.Name) == ".json" {
		var buf bytes.Buffer
		if err := json.Indent(&buf, file.Data, "", "  "); err != nil {
			return file, err
		}
		file.Data = buf.Bytes()
	}
	return file, nil
}
