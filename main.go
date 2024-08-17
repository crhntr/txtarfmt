package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"go/format"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/txtar"
)

func main() {
	var (
		formatGo   bool
		formatJSON bool
		ext        string
	)
	flag.BoolVar(&formatGo, "go", true, "format Go code")
	flag.BoolVar(&formatJSON, "json", true, "format JSON files")
	flag.StringVar(&ext, "ext", ".txtar", "file extension filter")
	flag.Parse()
	for _, arg := range flag.Args() {
		matches, err := filepath.Glob(arg)
		if err != nil {
			log.Fatal(err)
		}
		for _, match := range matches {
			if ext != "" && filepath.Ext(match) != ext {
				continue
			}
			archive, err := txtar.ParseFile(match)
			if err != nil {
				log.Fatal(err)
			}
			info, err := os.Stat(match)
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range archive.Files {
				file, err := formatFile(file, formatGo, formatJSON)
				if err != nil {
					log.Fatal(err)
				}
				if err := os.WriteFile(match, file.Data, info.Mode()); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func formatFile(file txtar.File, formatGo, formatJSON bool) (txtar.File, error) {
	if formatGo && filepath.Ext(file.Name) == ".go" {
		out, err := format.Source(file.Data)
		if err != nil {
			return file, err
		}
		file.Data = out
	}
	if formatJSON && filepath.Ext(file.Name) == ".json" {
		var buf bytes.Buffer
		if err := json.Indent(&buf, file.Data, "", "  "); err != nil {
			return file, err
		}
		file.Data = buf.Bytes()
	}
	return file, nil
}
