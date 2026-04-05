package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".emm"), nil
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
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", rel, err)
			}
			return nil
		}
		if _, err := os.Stat(target); err == nil {
			return nil
		}
		data, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", rel, err)
		}
		if err := os.WriteFile(target, data, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", rel, err)
		}
		return nil
	})
}
