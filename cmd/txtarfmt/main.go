package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/tools/txtar"

	"github.com/crhntr/txtarfmt"
)

func main() {
	value, ok := os.LookupEnv("TXTTARFMT_ALLOW_INSECURE")
	allowInsecure := ok && value == "true"
	if allowInsecure {
		fmt.Println("TXTTARFMT_ALLOW_INSECURE is true")
	}

	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	log.SetPrefix("txtarfmt: ")
	var (
		config txtarfmt.Configuration
		ext    string

		goGetU    bool
		goModTidy bool
	)
	flag.BoolVar(&config.SkipGo, "skip-go", false, "skip formatting Go code")
	flag.BoolVar(&config.SkipJSON, "skip-json", false, "skip formatting JSON files")
	flag.BoolVar(&config.SkipGoMod, "skip-gomod", false, "skip formatting go.mod files")
	flag.BoolVar(&config.SkipYAML, "skip-yaml", false, "skip formatting YAML files")
	if allowInsecure {
		flag.BoolVar(&goGetU, "INSECURE-go-get-u", false, "run go get -u on extracted archive")
		flag.BoolVar(&goModTidy, "INSECURE-go-mod-tidy", false, "run go mod tidy on extracted archive")
	}
	flag.StringVar(&ext, "ext", ".txtar", "file extension filter")
	flag.Parse()
	count := 0
	for i, arg := range flag.Args() {
		matches, err := filepath.Glob(arg)
		if err != nil {
			log.Fatalf("bad command argument %d: %s", i, err)
		}
		for _, match := range matches {
			if ext != "" && filepath.Ext(match) != ext {
				continue
			}
			log.Printf("glob match: %q\n", match)
			count++
			archive, err := txtar.ParseFile(match)
			if err != nil {
				log.Fatalf("%s: %s", match, err)
			}
			info, err := os.Stat(match)
			if err != nil {
				log.Fatalf("%s: %s", match, err)
			}

			if allowInsecure {
				var commands []*exec.Cmd
				if goGetU {
					commands = append(commands, exec.Command("go", "get", "-u"))
				}
				if goGetU || goModTidy {
					commands = append(commands, exec.Command("go", "mod", "tidy"))
				}
				if len(commands) > 0 {
					if err := execute(archive, commands); err != nil {
						log.Fatal(err)
					}
				}
			}

			if err := txtarfmt.Archive(archive, config); err != nil {
				log.Fatalf("%s: %s", match, err)
			}
			if err := os.WriteFile(match, txtar.Format(archive), info.Mode()); err != nil {
				log.Fatalf("%s: %s", match, err)
			}
		}
	}
	if count == 0 {
		log.Fatal("no files matched")
	}
}

func execute(a *txtar.Archive, commands []*exec.Cmd) error {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("creating temporary directory failed: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	updated, err := txtarfmt.Execute(dir, nil, a, commands...)
	if err != nil {
		return err
	}
	a.Files = updated.Files
	return nil
}
