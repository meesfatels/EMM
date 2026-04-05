package allowlist_test

import (
	"errors"
	"testing"

	"github.com/meesfatels/emm/internal/allowlist"
)

func TestEnforcer_ExactMatch(t *testing.T) {
	e := allowlist.NewEnforcer([]string{"ls", "cat", "git"})
	if err := e.Check("ls"); err != nil {
		t.Errorf("expected ls to be allowed, got %v", err)
	}
	if err := e.Check("ls -la"); err != nil {
		t.Errorf("expected 'ls -la' to be allowed (binary is ls), got %v", err)
	}
}

func TestEnforcer_WildcardMatch(t *testing.T) {
	e := allowlist.NewEnforcer([]string{"git*"})
	for _, cmd := range []string{"git", "git-clone", "git-config", "git status"} {
		if err := e.Check(cmd); err != nil {
			t.Errorf("expected %q to be allowed, got %v", cmd, err)
		}
	}
}

func TestEnforcer_NotAllowed(t *testing.T) {
	e := allowlist.NewEnforcer([]string{"ls"})
	err := e.Check("rm -rf /")
	if !errors.Is(err, allowlist.ErrNotAllowed) {
		t.Errorf("expected ErrNotAllowed, got %v", err)
	}
}

func TestEnforcer_EmptyCommand(t *testing.T) {
	e := allowlist.NewEnforcer([]string{"ls"})
	if err := e.Check(""); !errors.Is(err, allowlist.ErrEmptyCommand) {
		t.Errorf("expected ErrEmptyCommand, got %v", err)
	}
	if err := e.Check("   "); !errors.Is(err, allowlist.ErrEmptyCommand) {
		t.Errorf("expected ErrEmptyCommand for whitespace, got %v", err)
	}
}

func TestEnforcer_MultipleAllowlists(t *testing.T) {
	e := allowlist.NewEnforcer([]string{"ls", "cat"}, []string{"git*"})
	for _, cmd := range []string{"ls", "cat", "git", "git-clone"} {
		if err := e.Check(cmd); err != nil {
			t.Errorf("expected %q to be allowed, got %v", cmd, err)
		}
	}
	if err := e.Check("rm"); !errors.Is(err, allowlist.ErrNotAllowed) {
		t.Errorf("expected rm to be denied")
	}
}

func TestEnforcer_EmptyEnforcer(t *testing.T) {
	e := allowlist.NewEnforcer()
	if err := e.Check("ls"); !errors.Is(err, allowlist.ErrNotAllowed) {
		t.Errorf("expected everything to be denied with no allowlists")
	}
}
