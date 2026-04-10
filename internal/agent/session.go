package agent

import (
	"context"
	"encoding/json"
	"fmt"
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
	return &Session{
		agent:      a,
		minionName: minionName,
		username:   username,
		minion:     m,
		client:     client,
		messages:   []openrouter.Message{{Role: "system", Content: BuildPrompt(a)}},
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

func (s *Session) Send(ctx context.Context, content string, onToken func(string), onTool func(name, input, output string)) string {
	s.messages = append(s.messages, openrouter.Message{Role: "user", Content: content})

	var tools []openrouter.Tool
	toolMap := make(map[string]tool.Tool)
	for _, t := range s.agent.Tools {
		def := t.Definition()
		tools = append(tools, def)
		toolMap[def.Function.Name] = t
	}

	const maxRounds = 20
	for range maxRounds {
		req := &openrouter.Request{Minion: s.minion, Messages: s.messages, Tools: tools}
		stream := s.client.Complete(ctx, req)

		var resp strings.Builder
		for {
			token, ok := stream.Recv()
			if !ok {
				break
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
			s.messages = append(s.messages, openrouter.Message{Role: "assistant", Content: full})
			return full
		}

		s.messages = append(s.messages, openrouter.Message{Role: "assistant", ToolCalls: toolCalls})

		for _, tc := range toolCalls {
			t, ok := toolMap[tc.Function.Name]
			result := fmt.Sprintf("error: tool %s not found", tc.Function.Name)
			if ok {
				result = strings.TrimSpace(t.Execute(ctx, tc.Function.Arguments))
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

	return ""
}

func (s *Session) Save(emmDir, name string) {
	convsDir := filepath.Join(emmDir, "conversations")
	if err := os.MkdirAll(convsDir, 0o755); err != nil {
		panic(err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "---\nagent: %s\nminion: %s\n---\n\n", s.agent.Name, s.minionName)

	for _, msg := range s.messages {
		if msg.Role == "system" {
			continue
		}
		data, _ := json.Marshal(msg)
		fmt.Fprintf(&b, "<!-- message: %s -->\n", data)
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

	if err := os.WriteFile(filepath.Join(convsDir, name+".md"), []byte(b.String()), 0o644); err != nil {
		panic(err)
	}
}

// Load loads a conversation by name. Returns false if the file does not exist.
func (s *Session) Load(emmDir, name string) bool {
	data, err := os.ReadFile(filepath.Join(emmDir, "conversations", name+".md"))
	if err != nil {
		return false
	}

	var msgs []openrouter.Message
	if len(s.messages) > 0 && s.messages[0].Role == "system" {
		msgs = append(msgs, s.messages[0])
	}

	const prefix = "<!-- message: "
	const suffix = " -->"
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, prefix) || !strings.HasSuffix(line, suffix) {
			continue
		}
		var msg openrouter.Message
		if err := json.Unmarshal([]byte(line[len(prefix):len(line)-len(suffix)]), &msg); err != nil {
			continue
		}
		msgs = append(msgs, msg)
	}

	s.messages = msgs
	return true
}
