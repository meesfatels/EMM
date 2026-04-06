package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meesfatels/emm/internal/agent"
)

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "emm.yaml"), "api_key: sk-abc123\nusername: alice\n")

	c, err := agent.LoadConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.APIKey != "sk-abc123" {
		t.Errorf("APIKey: got %q, want sk-abc123", c.APIKey)
	}
	if c.Username != "alice" {
		t.Errorf("Username: got %q, want alice", c.Username)
	}
	if c.BaseURL != "https://openrouter.ai/api/v1/chat/completions" {
		t.Errorf("BaseURL: got %q, want default", c.BaseURL)
	}
}

func TestLoadConfig_MissingAPIKey(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "emm.yaml"), "username: bob\n")

	if _, err := agent.LoadConfig(dir); err == nil {
		t.Error("expected error for missing api_key")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "emm.yaml"), "api_key: sk-x\n")

	c, err := agent.LoadConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Username != "user" {
		t.Errorf("Username: got %q, want user", c.Username)
	}
	if c.BaseURL != "https://openrouter.ai/api/v1/chat/completions" {
		t.Errorf("BaseURL: got %q, want default", c.BaseURL)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
