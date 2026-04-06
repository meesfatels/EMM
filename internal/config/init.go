package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func fileModeFor(rel string) os.FileMode {
	if rel == "emm.yaml" {
		return 0o600
	}
	return 0o644
}

func Init(templateFS fs.FS) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	return fs.WalkDir(templateFS, ".EMM", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking template: %w", err)
		}
		rel, err := filepath.Rel(".EMM", path)
		if err != nil {
			return fmt.Errorf("resolving path %s: %w", path, err)
		}
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if _, err := os.Stat(target); err == nil {
			return nil
		}
		data, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", rel, err)
		}
		return os.WriteFile(target, data, fileModeFor(rel))
	})
}
