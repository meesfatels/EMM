package runtime

import (
	"fmt"
	"strings"

	"github.com/meesfatels/emm/internal/loader"
)

func BuildPrompt(instinct *loader.Instinct) string {
	var b strings.Builder
	for _, f := range instinct.Files {
		content, ok := instinct.Content[f.Name]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "[%s: %s]\n%s\n\n", f.Name, f.Interpretation, strings.TrimSpace(content))
	}
	return strings.TrimSpace(b.String())
}
