package tool

import (
	"context"

	"github.com/meesfatels/emm/internal/openrouter"
)

// Tool defines the interface for model-callable tools.
type Tool interface {
	// Definition returns the OpenRouter tool specification.
	Definition() openrouter.Tool
	// Execute runs the tool with the provided JSON arguments (as a string).
	Execute(ctx context.Context, args string) (string, error)
}
