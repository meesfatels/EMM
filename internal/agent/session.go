package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/meesfatels/emm/internal/minion"
	"github.com/meesfatels/emm/internal/openrouter"
	"github.com/meesfatels/emm/internal/shell"
)

// runShellTool is the OpenRouter tool definition for shell execution.
var runShellTool = openrouter.Tool{
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

type Session struct {
	agent      *Agent
	minionName string
	username   string
	minion     minion.Minion
	client     *openrouter.Client
	executor   *shell.Executor
	messages   []openrouter.Message
}

func NewSession(a *Agent, minionName string, m minion.Minion, client *openrouter.Client, username string) *Session {
	prompt := BuildPrompt(a)
	return &Session{
		agent:      a,
		minionName: minionName,
		username:   username,
		minion:     m,
		client:     client,
		executor:   shell.NewExecutor(a.Shell),
		messages: []openrouter.Message{
			{Role: "system", Content: prompt},
		},
	}
}

func (s *Session) SwitchAgent(a *Agent) {
	s.agent = a
	s.executor = shell.NewExecutor(a.Shell)
}

func (s *Session) SwitchMinion(m minion.Minion, name string) {
	s.minion = m
	s.minionName = name
}

func (s *Session) Messages() []openrouter.Message {
	return s.messages
}

// Send sends content to the model, streaming tokens via onToken.
// If the agent has shell rules, tool calls are handled transparently:
// onShell is called for each execution with the command and its output.
// The session retries until the model produces a plain text response.
func (s *Session) Send(ctx context.Context, content string, onToken func(string), onShell func(cmd, output string)) (string, error) {
	s.messages = append(s.messages, openrouter.Message{
		Role:    "user",
		Content: content,
	})
	startLen := len(s.messages)

	var tools []openrouter.Tool
	if len(s.agent.Shell) > 0 {
		tools = []openrouter.Tool{runShellTool}
	}

	const maxRounds = 20
	for round := 0; round < maxRounds; round++ {
		req := openrouter.NewRequest(s.minion, s.messages, tools)
		stream, err := s.client.Complete(ctx, req)
		if err != nil {
			s.messages = s.messages[:startLen-1]
			return "", fmt.Errorf("completing request: %w", err)
		}

		var resp strings.Builder
		for {
			token, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				stream.Close()
				s.messages = s.messages[:startLen-1]
				return resp.String(), fmt.Errorf("streaming response: %w", err)
			}
			resp.WriteString(token)
			if onToken != nil {
				onToken(token)
			}
		}

		toolCalls := stream.ToolCalls()
		stream.Close()

		if len(toolCalls) == 0 {
			full := resp.String()
			s.messages = append(s.messages, openrouter.Message{
				Role:    "assistant",
				Content: full,
			})
			return full, nil
		}

		// Append the assistant message that contains the tool calls.
		s.messages = append(s.messages, openrouter.Message{
			Role:      "assistant",
			ToolCalls: toolCalls,
		})

		// Execute each tool call and append the results.
		for _, tc := range toolCalls {
			var args struct {
				Cmd string `json:"cmd"`
			}
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				s.messages = s.messages[:startLen-1]
				return "", fmt.Errorf("parsing tool arguments: %w", err)
			}

			output, runErr := s.executor.Run(ctx, args.Cmd)
			result := strings.TrimSpace(output)
			if runErr != nil {
				if result != "" {
					result += "\n[exit error: " + runErr.Error() + "]"
				} else {
					result = "[exit error: " + runErr.Error() + "]"
				}
			}
			if result == "" {
				result = "(no output)"
			}

			if onShell != nil {
				onShell(args.Cmd, result)
			}

			s.messages = append(s.messages, openrouter.Message{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    result,
			})
		}
		// Loop to get the model's response after tool execution.
	}
	s.messages = s.messages[:startLen-1]
	return "", fmt.Errorf("exceeded maximum tool call rounds (%d)", maxRounds)
}

func (s *Session) Save(emmDir, name string) error {
	convsDir := filepath.Join(emmDir, "conversations")
	if err := os.MkdirAll(convsDir, 0o755); err != nil {
		return fmt.Errorf("creating conversations dir: %w", err)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# %s — %s/%s\n\n", name, s.agent.Name, s.minionName)
	for _, msg := range s.messages {
		switch msg.Role {
		case "user":
			fmt.Fprintf(&b, "## %s\n\n%s\n\n", s.username, msg.Content)
		case "assistant":
			if msg.Content != "" {
				fmt.Fprintf(&b, "## %s-%s\n\n%s\n\n", s.agent.Name, s.minionName, msg.Content)
			}
		}
	}
	return os.WriteFile(filepath.Join(convsDir, name+".md"), []byte(b.String()), 0o644)
}

func (s *Session) Load(emmDir, name string) error {
	convPath := filepath.Join(emmDir, "conversations", name+".md")
	data, err := os.ReadFile(convPath)
	if err != nil {
		return fmt.Errorf("reading conversation file: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	var newMessages []openrouter.Message
	if len(s.messages) > 0 && s.messages[0].Role == "system" {
		newMessages = append(newMessages, s.messages[0])
	}

	var currentRole string
	var currentContent strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if currentRole != "" {
				newMessages = append(newMessages, openrouter.Message{
					Role:    currentRole,
					Content: strings.TrimSpace(currentContent.String()),
				})
				currentContent.Reset()
			}
			rolePart := strings.TrimPrefix(line, "## ")
			if rolePart == s.username {
				currentRole = "user"
			} else {
				currentRole = "assistant"
			}
			continue
		}
		if currentRole != "" {
			currentContent.WriteString(line + "\n")
		}
	}
	if currentRole != "" {
		newMessages = append(newMessages, openrouter.Message{
			Role:    currentRole,
			Content: strings.TrimSpace(currentContent.String()),
		})
	}

	s.messages = newMessages
	return nil
}
