package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meesfatels/emm/internal/agent"
)

// ---- helpers ----------------------------------------------------------------

func panics(fn func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// makeAgent creates a minimal valid agent directory and returns the emm root dir.
func makeAgent(t *testing.T, name string, instinctFiles map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	agentDir := filepath.Join(dir, "agents", name)

	var igYAML string
	for fname := range instinctFiles {
		igYAML += "  - name: " + fname + "\n    interpretation: test\n"
		writeFile(t, filepath.Join(agentDir, "instinct", fname), instinctFiles[fname])
	}
	writeFile(t, filepath.Join(agentDir, "instinct_guide.yaml"), "instinct:\n"+igYAML)
	return dir
}

// ---- BuildPrompt ------------------------------------------------------------

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name     string
		a        *agent.Agent
		contains []string
		empty    bool
	}{
		{
			name:  "empty agent",
			a:     &agent.Agent{},
			empty: true,
		},
		{
			name: "single instinct file",
			a: &agent.Agent{
				Instinct: []agent.InstinctFile{{Name: "p.md", Interpretation: "Personality"}},
				Content:  map[string]string{"p.md": "Be helpful."},
			},
			contains: []string{"[p.md: Personality]", "Be helpful."},
		},
		{
			name: "missing instinct file skipped",
			a: &agent.Agent{
				Instinct: []agent.InstinctFile{
					{Name: "missing.md", Interpretation: "Gone"},
					{Name: "present.md", Interpretation: "Here"},
				},
				Content: map[string]string{"present.md": "I exist."},
			},
			contains: []string{"I exist."},
		},
		{
			name: "multiple files joined",
			a: &agent.Agent{
				Instinct: []agent.InstinctFile{
					{Name: "a.md", Interpretation: "A"},
					{Name: "b.md", Interpretation: "B"},
				},
				Content: map[string]string{"a.md": "Content A", "b.md": "Content B"},
			},
			contains: []string{"Content A", "Content B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agent.BuildPrompt(tt.a)
			if tt.empty && got != "" {
				t.Errorf("expected empty prompt, got %q", got)
			}
			for _, want := range tt.contains {
				if !contains(got, want) {
					t.Errorf("prompt missing %q\ngot: %q", want, got)
				}
			}
		})
	}
}

// ---- LoadConfig -------------------------------------------------------------

func TestLoadConfig(t *testing.T) {
	t.Run("valid with defaults", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "emm.yaml"), "api_key: sk-test\n")
		cfg := agent.LoadConfig(dir)
		if cfg.APIKey != "sk-test" {
			t.Errorf("APIKey: got %q", cfg.APIKey)
		}
		if cfg.Username != "user" {
			t.Errorf("Username default: got %q", cfg.Username)
		}
		if cfg.DefaultAgent != "example" {
			t.Errorf("DefaultAgent default: got %q", cfg.DefaultAgent)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "emm.yaml"),
			"api_key: sk-abc\nusername: alice\ndefault_agent: myagent\ndefault_minion: myminion\n")
		cfg := agent.LoadConfig(dir)
		if cfg.Username != "alice" {
			t.Errorf("Username: got %q", cfg.Username)
		}
		if cfg.DefaultAgent != "myagent" {
			t.Errorf("DefaultAgent: got %q", cfg.DefaultAgent)
		}
	})

	t.Run("missing api_key panics", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "emm.yaml"), "username: bob\n")
		if !panics(func() { agent.LoadConfig(dir) }) {
			t.Error("expected panic for missing api_key")
		}
	})

	t.Run("missing file panics", func(t *testing.T) {
		if !panics(func() { agent.LoadConfig(t.TempDir()) }) {
			t.Error("expected panic for missing emm.yaml")
		}
	})
}

// ---- Load -------------------------------------------------------------------

func TestLoad(t *testing.T) {
	t.Run("loads instinct files", func(t *testing.T) {
		dir := makeAgent(t, "myagent", map[string]string{
			"personality.md": "Be helpful.",
			"rules.md":       "Follow rules.",
		})
		a := agent.Load(dir, "myagent")
		if a.Name != "myagent" {
			t.Errorf("Name: got %q", a.Name)
		}
		if len(a.Instinct) != 2 {
			t.Errorf("Instinct count: got %d, want 2", len(a.Instinct))
		}
		if a.Content["personality.md"] != "Be helpful." {
			t.Errorf("Content: got %q", a.Content["personality.md"])
		}
	})

	t.Run("missing instinct_guide panics", func(t *testing.T) {
		dir := t.TempDir()
		if !panics(func() { agent.Load(dir, "nonexistent") }) {
			t.Error("expected panic for missing agent")
		}
	})

	t.Run("allowlist optional", func(t *testing.T) {
		dir := makeAgent(t, "noallowlist", map[string]string{"p.md": "hi"})
		a := agent.Load(dir, "noallowlist")
		if len(a.Tools) != 0 {
			t.Errorf("expected no tools without allowlist, got %d", len(a.Tools))
		}
	})

	t.Run("allowlist adds shell tool", func(t *testing.T) {
		dir := makeAgent(t, "withallowlist", map[string]string{"p.md": "hi"})
		allowlistPath := filepath.Join(dir, "agents", "withallowlist", "allowlist.yaml")
		writeFile(t, allowlistPath, "shell:\n  - binary: echo\n")
		a := agent.Load(dir, "withallowlist")
		if len(a.Tools) != 1 {
			t.Errorf("expected 1 tool with allowlist, got %d", len(a.Tools))
		}
	})
}

func TestLoadAll(t *testing.T) {
	dir := t.TempDir()
	makeAgentInDir(t, dir, "alpha", map[string]string{"a.md": "A"})
	makeAgentInDir(t, dir, "beta", map[string]string{"b.md": "B"})

	agents := agent.LoadAll(dir)
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}
	if _, ok := agents["alpha"]; !ok {
		t.Error("missing agent alpha")
	}
	if _, ok := agents["beta"]; !ok {
		t.Error("missing agent beta")
	}
}

// ---- helpers ----------------------------------------------------------------

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}

// makeAgentInDir creates an agent inside an already-existing root dir.
func makeAgentInDir(t *testing.T, dir, name string, instinctFiles map[string]string) {
	t.Helper()
	agentDir := filepath.Join(dir, "agents", name)
	var igYAML string
	for fname := range instinctFiles {
		igYAML += "  - name: " + fname + "\n    interpretation: test\n"
		writeFile(t, filepath.Join(agentDir, "instinct", fname), instinctFiles[fname])
	}
	writeFile(t, filepath.Join(agentDir, "instinct_guide.yaml"), "instinct:\n"+igYAML)
}
