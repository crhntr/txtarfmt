package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/txtar"

	"github.com/crhntr/txtarfmt"
)

func main() {
	var (
		config txtarfmt.Configuration
		ext    string
	)
	flag.BoolVar(&config.SkipGo, "skip-go", false, "skip formatting Go code")
	flag.BoolVar(&config.SkipJSON, "skip-json", false, "skip formatting JSON files")
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
			if err := txtarfmt.Archive(archive, config); err != nil {
				log.Fatal(err)
			}
			if err := os.WriteFile(match, txtar.Format(archive), info.Mode()); err != nil {
				log.Fatal(err)
			}
		}
	}
}
