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
	if len(s.messages) > 0 && s.messages[0].Role == "system" {
		s.messages[0].Content = BuildPrompt(a)
	}
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
		if msg.Role == "system" {
			continue
		}
		// Machine-readable comment for lossless round-trip (tool calls, IDs, etc.).
		// json.Marshal HTML-escapes '<', '>', '&' so the output never contains "-->".
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("encoding message: %w", err)
		}
		fmt.Fprintf(&b, "<!-- message: %s -->\n", data)
		// Human-readable section below the comment.
		switch msg.Role {
		case "user":
			fmt.Fprintf(&b, "## %s\n\n%s\n\n", s.username, strings.TrimSpace(msg.Content))
		case "assistant":
			if msg.Content != "" {
				fmt.Fprintf(&b, "## %s\n\n%s\n\n", s.agent.Name, strings.TrimSpace(msg.Content))
			} else {
				fmt.Fprintf(&b, "## %s\n\n[%d tool call(s)]\n\n", s.agent.Name, len(msg.ToolCalls))
			}
		case "tool":
			fmt.Fprintf(&b, "## tool\n\n%s\n\n", strings.TrimSpace(msg.Content))
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

	var newMessages []openrouter.Message
	if len(s.messages) > 0 && s.messages[0].Role == "system" {
		newMessages = append(newMessages, s.messages[0])
	}

	content := string(data)
	const msgPrefix = "<!-- message: "
	const msgSuffix = " -->"

	if strings.Contains(content, msgPrefix) {
		// Current format: one JSON-encoded message per <!-- message: {...} --> line.
		for _, line := range strings.Split(content, "\n") {
			if !strings.HasPrefix(line, msgPrefix) || !strings.HasSuffix(line, msgSuffix) {
				continue
			}
			jsonStr := line[len(msgPrefix) : len(line)-len(msgSuffix)]
			var msg openrouter.Message
			if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
				continue
			}
			newMessages = append(newMessages, msg)
		}
	} else {
		// Legacy format: <!-- role: X --> ... sections with text content only.
		for _, sec := range strings.Split(content, "<!-- role: ") {
			if !strings.Contains(sec, " -->") {
				continue
			}
			parts := strings.SplitN(sec, " -->", 2)
			role := strings.TrimSpace(parts[0])
			body := parts[1]
			msgParts := strings.SplitN(body, "\n\n", 2)
			if len(msgParts) < 2 {
				continue
			}
			newMessages = append(newMessages, openrouter.Message{
				Role:    role,
				Content: strings.TrimSpace(msgParts[1]),
			})
		}
	}

	s.messages = newMessages
	return nil
}
