package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestInit_UsesSecureModeForConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	tpl := fstest.MapFS{
		".EMM/emm.yaml":       &fstest.MapFile{Data: []byte("api_key: test\n"), Mode: fs.FileMode(0o644)},
		".EMM/minions/x.yaml": &fstest.MapFile{Data: []byte("model: test\n"), Mode: fs.FileMode(0o644)},
	}

	if err := Init(tpl); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	emmPath := filepath.Join(home, ".emm", "emm.yaml")
	st, err := os.Stat(emmPath)
	if err != nil {
		t.Fatalf("stat emm.yaml: %v", err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Fatalf("emm.yaml perms = %o, want 600", st.Mode().Perm())
	}

}
