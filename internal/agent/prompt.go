package agent

import (
	"fmt"
	"strings"
)

func BuildPrompt(a *Agent) string {
	var b strings.Builder
	for _, f := range a.Instinct {
		content, ok := a.Content[f.Name]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "[%s: %s]\n%s\n\n", f.Name, f.Interpretation, strings.TrimSpace(content))
	}
	return strings.TrimSpace(b.String())
}
