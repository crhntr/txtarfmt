package txtarfmt

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"

	"golang.org/x/tools/txtar"
)

func Execute(tmpDir string, ignoreGlobs []string, archive *txtar.Archive, commands ...*exec.Cmd) (*txtar.Archive, error) {
	tmpRoot, err := os.OpenRoot(tmpDir)
	if err != nil {
		return nil, err
	}

	if archive != nil {
		for _, file := range archive.Files {
			dir := filepath.Dir(filepath.FromSlash(file.Name))
			if err := tmpRoot.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("copy archive directory failed: %w", err)
			}
			if err := tmpRoot.WriteFile(filepath.FromSlash(file.Name), file.Data, 0644); err != nil {
				return nil, fmt.Errorf("copy archive file failed: %w", err)
			}
		}
	}

	for i, cmd := range commands {
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("command %d failed: %w", i, err)
		}
	}

	tmpFS := tmpRoot.FS()

	var filesPaths []string
	if err := fs.WalkDir(tmpFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return err
		}
		filesPaths = append(filesPaths, path)
		return err
	}); err != nil {
		return nil, err
	}

	filesPaths = slices.DeleteFunc(filesPaths, func(s string) bool {
		for _, g := range ignoreGlobs {
			if ok, err := path.Match(g, s); err == nil && ok {
				return true
			}
		}
		return false
	})

	slices.Sort(filesPaths)

	result := new(txtar.Archive)
	if archive != nil {
		result.Comment = archive.Comment
	}

	for _, p := range filesPaths {
		buf, err := fs.ReadFile(tmpFS, p)
		if err != nil {
			return nil, err
		}
		result.Files = append(result.Files, txtar.File{
			Name: p,
			Data: buf,
		})
	}

	return result, nil
}
