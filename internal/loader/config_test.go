package loader_test

import (
	"testing"

	"github.com/meesfatels/emm/internal/loader"
)

func TestConfig_APIKey(t *testing.T) {
	t.Run("valid key", func(t *testing.T) {
		c := loader.Config{"api_key": "sk-abc123"}
		key, err := c.APIKey()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key != "sk-abc123" {
			t.Errorf("expected sk-abc123, got %q", key)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		c := loader.Config{}
		if _, err := c.APIKey(); err == nil {
			t.Error("expected error for missing api_key")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		c := loader.Config{"api_key": ""}
		if _, err := c.APIKey(); err == nil {
			t.Error("expected error for empty api_key")
		}
	})

	t.Run("non-string key", func(t *testing.T) {
		c := loader.Config{"api_key": 42}
		if _, err := c.APIKey(); err == nil {
			t.Error("expected error for non-string api_key")
		}
	})
}

func TestConfig_BaseURL(t *testing.T) {
	const defaultURL = "https://openrouter.ai/api/v1/chat/completions"

	t.Run("default when missing", func(t *testing.T) {
		c := loader.Config{}
		if got := c.BaseURL(); got != defaultURL {
			t.Errorf("expected default URL, got %q", got)
		}
	})

	t.Run("default when empty", func(t *testing.T) {
		c := loader.Config{"base_url": ""}
		if got := c.BaseURL(); got != defaultURL {
			t.Errorf("expected default URL for empty value, got %q", got)
		}
	})

	t.Run("custom URL", func(t *testing.T) {
		c := loader.Config{"base_url": "https://custom.example.com/v1"}
		if got := c.BaseURL(); got != "https://custom.example.com/v1" {
			t.Errorf("expected custom URL, got %q", got)
		}
	})
}
