package allowlist
import (
	"errors"
	"fmt"
	"strings"
)
var (
	ErrEmptyCommand = errors.New("empty command")
	ErrNotAllowed   = errors.New("command not allowed")
)
type Enforcer struct {
	patterns []string
}
func NewEnforcer(allowlists ...[]string) *Enforcer {
	var total int
	for _, al := range allowlists {
		total += len(al)
	}
	patterns := make([]string, 0, total)
	for _, al := range allowlists {
		patterns = append(patterns, al...)
	}
	return &Enforcer{patterns: patterns}
}
func (e *Enforcer) Check(command string) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ErrEmptyCommand
	}
	bin := fields[0]
	for _, p := range e.patterns {
		if strings.HasSuffix(p, "*") {
			if strings.HasPrefix(bin, strings.TrimSuffix(p, "*")) {
				return nil
			}
		} else if bin == p {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrNotAllowed, bin)
}
