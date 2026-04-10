package shell_test

import (
	"context"
	"strings"
	"testing"

	"github.com/meesfatels/emm/internal/shell"
)

var (
	ctx   = context.Background()
	rules = []shell.Rule{
		{Binary: "echo"},
		{Binary: "true"},
		{Binary: "git", Deny: []string{"push", "reset --hard"}},
	}
)

func executor(r ...shell.Rule) *shell.Executor {
	return shell.NewExecutor(r)
}

func TestRun(t *testing.T) {
	e := executor(rules...)

	tests := []struct {
		name    string
		cmd     string
		wantOut string // exact, or empty to just check prefix
		wantErr bool   // true if output should start with "error:"
	}{
		{"runs allowed command", "echo hello", "hello\n", false},
		{"empty output on success", "true", "", false},
		{"blocks unknown binary", "ls -la", "", true},
		{"blocks denied subcommand", "git push origin main", "", true},
		{"blocks another denied subcommand", "git reset --hard HEAD", "", true},
		{"allows non-denied git subcommand", "git --version", "", false},
		{"blocks malformed command", "echo 'unclosed", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := e.Run(ctx, tt.cmd)
			if tt.wantErr {
				if !strings.HasPrefix(out, "error:") {
					t.Errorf("Run(%q) = %q, want error: prefix", tt.cmd, out)
				}
				return
			}
			if tt.wantOut != "" && out != tt.wantOut {
				t.Errorf("Run(%q) = %q, want %q", tt.cmd, out, tt.wantOut)
			}
			if strings.HasPrefix(out, "error:") {
				t.Errorf("Run(%q) unexpected error: %q", tt.cmd, out)
			}
		})
	}
}

func TestRun_EmptyRules(t *testing.T) {
	e := executor()
	out := e.Run(ctx, "echo hi")
	if !strings.HasPrefix(out, "error:") {
		t.Errorf("expected error with empty rules, got %q", out)
	}
}

func TestExecute(t *testing.T) {
	e := executor(shell.Rule{Binary: "echo"})

	t.Run("valid JSON args", func(t *testing.T) {
		out := e.Execute(ctx, `{"cmd":"echo hi"}`)
		if out != "hi\n" {
			t.Errorf("got %q, want %q", out, "hi\n")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		out := e.Execute(ctx, "not json")
		if !strings.HasPrefix(out, "error:") {
			t.Errorf("expected error, got %q", out)
		}
	})

	t.Run("missing cmd field", func(t *testing.T) {
		out := e.Execute(ctx, `{"other":"field"}`)
		// empty cmd string → shlex produces empty parts → not allowed
		if !strings.HasPrefix(out, "error:") {
			t.Errorf("expected error, got %q", out)
		}
	})
}
