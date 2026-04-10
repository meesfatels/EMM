package tool

import (
	"context"

	"github.com/meesfatels/emm/internal/openrouter"
)

// Tool defines the interface for model-callable tools.
type Tool interface {
	Definition() openrouter.Tool
	Execute(ctx context.Context, args string) string
}
