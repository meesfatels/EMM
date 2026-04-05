package runtime

import (
	"github.com/meesfatels/emm/internal/loader"
	"strings"
)

func BuildPrompt(instinct *loader.Instinct) string {
	var b strings.Builder
	for _, f := range instinct.Guide.Files {
		content, ok := instinct.Content[f.Name]
		if !ok {
			continue
		}
		b.WriteString("[")
		b.WriteString(f.Name)
		b.WriteString(": ")
		b.WriteString(f.Interpretation)
		b.WriteString("]\n")
		b.WriteString(strings.TrimSpace(content))
		b.WriteString("\n\n")
	}
	return strings.TrimSpace(b.String())
}
