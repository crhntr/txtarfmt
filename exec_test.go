package txtarfmt_test

import (
	"bytes"
	"os/exec"
	"slices"
	"testing"

	"golang.org/x/tools/txtar"

	"github.com/crhntr/txtarfmt"
)

func TestExecute(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		dir := t.TempDir()

		result, err := txtarfmt.Execute(dir, nil, nil)
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		if result == nil {
			t.Error("expected non-nil result")
		}
	})

	t.Run("go mod init", func(t *testing.T) {
		dir := t.TempDir()

		result, err := txtarfmt.Execute(dir, nil, nil, exec.CommandContext(t.Context(), "go", "mod", "init", "example.com"))
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}

		i := slices.IndexFunc(result.Files, func(f txtar.File) bool {
			return f.Name == "go.mod"
		})
		if i < 0 {
			t.Fatal("expected go mod file")
		}
		if !bytes.Contains(result.Files[i].Data, []byte("module example.com")) {
			t.Log(string(result.Files[i].Data))
			t.Errorf("expected to find example.com in module file")
		}
	})

	t.Run("go fmt", func(t *testing.T) {
		dir := t.TempDir()

		result, err := txtarfmt.Execute(dir, nil, &txtar.Archive{
			Files: []txtar.File{
				{Name: "go.mod", Data: []byte("module example.com")},
				{Name: "cmd/hello/main.go", Data: []byte(`package main

func main(){println(     "hello"        ) }`)},
			},
		}, exec.CommandContext(t.Context(), "go", "fmt", "./..."))
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}

		i := slices.IndexFunc(result.Files, func(f txtar.File) bool {
			t.Log(f.Name)
			return f.Name == "cmd/hello/main.go"
		})
		if i < 0 {
			t.Fatal("expected file not found")
		}
		if !bytes.Contains(result.Files[i].Data, []byte(`func main() { println("hello") }`)) {
			t.Log(string(result.Files[i].Data))
			t.Errorf("expected to find example.com in module file")
		}
	})
}
