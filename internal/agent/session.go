package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/meesfatels/emm/internal/minion"
	"github.com/meesfatels/emm/internal/openrouter"
	"github.com/meesfatels/emm/internal/tool"
)

type Session struct {
	agent      *Agent
	minionName string
	username   string
	minion     minion.Minion
	client     *openrouter.Client
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
		messages: []openrouter.Message{
			{Role: "system", Content: prompt},
		},
	}
}

func (s *Session) SwitchAgent(a *Agent) {
	s.agent = a
}

func (s *Session) SwitchMinion(m minion.Minion, name string) {
	s.minion = m
	s.minionName = name
}

func (s *Session) Messages() []openrouter.Message {
	return s.messages
}

// Send sends content to the model, streaming tokens via onToken.
// Tool calls are handled transparently using the agent's tools.
// onTool is called for each execution with the tool name, input, and output.
// The session retries until the model produces a plain text response.
func (s *Session) Send(ctx context.Context, content string, onToken func(string), onTool func(name, input, output string)) (string, error) {
	s.messages = append(s.messages, openrouter.Message{
		Role:    "user",
		Content: content,
	})
	startLen := len(s.messages)

	var tools []openrouter.Tool
	toolMap := make(map[string]tool.Tool)
	for _, t := range s.agent.Tools {
		def := t.Definition()
		tools = append(tools, def)
		toolMap[def.Function.Name] = t
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
			t, ok := toolMap[tc.Function.Name]
			if !ok {
				s.messages = append(s.messages, openrouter.Message{
					Role:       "tool",
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf("error: tool %s not found", tc.Function.Name),
				})
				continue
			}

			output, runErr := t.Execute(ctx, tc.Function.Arguments)
			result := strings.TrimSpace(output)
			if runErr != nil {
				if result != "" {
					result += "\n[error: " + runErr.Error() + "]"
				} else {
					result = "[error: " + runErr.Error() + "]"
				}
			}
			if result == "" {
				result = "(no output)"
			}

			if onTool != nil {
				onTool(tc.Function.Name, tc.Function.Arguments, result)
			}

			s.messages = append(s.messages, openrouter.Message{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    result,
			})
		}
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
	fmt.Fprintf(&b, "---\nagent: %s\nminion: %s\n---\n\n", s.agent.Name, s.minionName)

	for _, msg := range s.messages {
		if msg.Content == "" {
			continue
		}
		if msg.Role != "user" && msg.Role != "assistant" {
			continue
		}

		displayName := s.username
		if msg.Role == "assistant" {
			displayName = s.agent.Name
		}

		fmt.Fprintf(&b, "<!-- role: %s -->\n## %s\n\n%s\n\n", msg.Role, displayName, strings.TrimSpace(msg.Content))
	}
	return os.WriteFile(filepath.Join(convsDir, name+".md"), []byte(b.String()), 0o644)
}

func (s *Session) Load(emmDir, name string) error {
	convPath := filepath.Join(emmDir, "conversations", name+".md")
	data, err := os.ReadFile(convPath)
	if err != nil {
		return fmt.Errorf("reading conversation file: %w", err)
	}

	content := string(data)
	sections := strings.Split(content, "<!-- role: ")

	var newMessages []openrouter.Message
	if len(s.messages) > 0 && s.messages[0].Role == "system" {
		newMessages = append(newMessages, s.messages[0])
	}

	for _, sec := range sections {
		if !strings.Contains(sec, " -->") {
			continue
		}
		parts := strings.SplitN(sec, " -->", 2)
		role := strings.TrimSpace(parts[0])
		body := parts[1]

		// Find the first double newline after the header to get the actual content
		// The header looks like "\n## Name\n\n"
		msgParts := strings.SplitN(body, "\n\n", 2)
		if len(msgParts) < 2 {
			continue
		}
		msgContent := strings.TrimSpace(msgParts[1])

		newMessages = append(newMessages, openrouter.Message{
			Role:    role,
			Content: msgContent,
		})
	}

	s.messages = newMessages
	return nil
}
