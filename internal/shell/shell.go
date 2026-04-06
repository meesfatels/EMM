package shell

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
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
