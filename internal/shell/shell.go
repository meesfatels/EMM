package shell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/meesfatels/emm/internal/openrouter"
)

type Rule struct {
	Binary string   `yaml:"binary"`
	Deny   []string `yaml:"deny"`
}

type Executor struct {
	rules []Rule
}

func NewExecutor(rules []Rule) *Executor {
	return &Executor{rules: rules}
}

func (e *Executor) Definition() openrouter.Tool {
	return openrouter.Tool{
		Type: "function",
		Function: openrouter.ToolDefinition{
			Name:        "run_shell",
			Description: "Execute a shell command and return its output.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"cmd": map[string]any{
						"type":        "string",
						"description": "The shell command to execute.",
					},
				},
				"required": []string{"cmd"},
			},
		},
	}
}

func (e *Executor) Execute(ctx context.Context, args string) (string, error) {
	var a struct {
		Cmd string `json:"cmd"`
	}
	if err := json.Unmarshal([]byte(args), &a); err != nil {
		return "", fmt.Errorf("parsing tool arguments: %w", err)
	}
	return e.Run(ctx, a.Cmd)
}

// Allowed reports whether cmd is permitted by the allowlist.
func (e *Executor) Allowed(cmd string) bool {
	if len(e.rules) == 0 {
		return false
	}
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}
	binary := parts[0]
	args := strings.Join(parts[1:], " ")
	for _, rule := range e.rules {
		if rule.Binary != binary {
			continue
		}
		for _, deny := range rule.Deny {
			if strings.Contains(args, deny) {
				return false
			}
		}
		return true
	}
	return false
}

// Run executes cmd if allowed, returning combined stdout+stderr.
// A non-zero exit code is returned as an error alongside any output.
func (e *Executor) Run(ctx context.Context, cmd string) (string, error) {
	if !e.Allowed(cmd) {
		return "", fmt.Errorf("command not allowed: %s", cmd)
	}
	parts := strings.Fields(cmd)
	c := exec.CommandContext(ctx, parts[0], parts[1:]...)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	return out.String(), err
}
