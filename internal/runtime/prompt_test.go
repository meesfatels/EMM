package runtime_test

import (
	"strings"
	"testing"

	"github.com/meesfatels/emm/internal/loader"
	"github.com/meesfatels/emm/internal/runtime"
)

func TestBuildPrompt_SingleFile(t *testing.T) {
	instinct := &loader.Instinct{
		Guide: loader.InstinctGuide{
			Files: []loader.InstinctFile{
				{Name: "personality.md", Interpretation: "Core personality"},
			},
		},
		Content: map[string]string{
			"personality.md": "You are a helpful assistant.",
		},
	}
	got := runtime.BuildPrompt(instinct)
	if !strings.Contains(got, "[personality.md: Core personality]") {
		t.Errorf("expected header in prompt, got: %q", got)
	}
	if !strings.Contains(got, "You are a helpful assistant.") {
		t.Errorf("expected content in prompt, got: %q", got)
	}
}

func TestBuildPrompt_MultipleFiles(t *testing.T) {
	instinct := &loader.Instinct{
		Guide: loader.InstinctGuide{
			Files: []loader.InstinctFile{
				{Name: "a.md", Interpretation: "First"},
				{Name: "b.md", Interpretation: "Second"},
			},
		},
		Content: map[string]string{
			"a.md": "Content A",
			"b.md": "Content B",
		},
	}
	got := runtime.BuildPrompt(instinct)
	if !strings.Contains(got, "Content A") || !strings.Contains(got, "Content B") {
		t.Errorf("expected both files in prompt, got: %q", got)
	}
	if strings.HasSuffix(got, "\n\n") {
		t.Errorf("prompt should be trimmed, got trailing whitespace")
	}
}

func TestBuildPrompt_MissingFile(t *testing.T) {
	instinct := &loader.Instinct{
		Guide: loader.InstinctGuide{
			Files: []loader.InstinctFile{
				{Name: "missing.md", Interpretation: "Won't be found"},
				{Name: "present.md", Interpretation: "Present"},
			},
		},
		Content: map[string]string{
			"present.md": "I am here",
		},
	}
	got := runtime.BuildPrompt(instinct)
	if strings.Contains(got, "missing.md") {
		t.Errorf("missing file should be skipped, got: %q", got)
	}
	if !strings.Contains(got, "I am here") {
		t.Errorf("present file should be included, got: %q", got)
	}
}

func TestBuildPrompt_Empty(t *testing.T) {
	instinct := &loader.Instinct{
		Guide:   loader.InstinctGuide{},
		Content: map[string]string{},
	}
	got := runtime.BuildPrompt(instinct)
	if got != "" {
		t.Errorf("expected empty prompt, got: %q", got)
	}
}
