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

	if len(a.Tools) > 0 {
		fmt.Fprintf(&b, "### Available Tools\n")
		for _, t := range a.Tools {
			def := t.Definition()
			fmt.Fprintf(&b, "- %s: %s\n", def.Function.Name, def.Function.Description)
		}
		fmt.Fprintf(&b, "\nYou have access to the tools above. Use them when needed to fulfill the user's request.\n")
	}

	return strings.TrimSpace(b.String())
}
