package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/meesfatels/emm/internal/loader"
	"github.com/meesfatels/emm/internal/openrouter"
)

type Session struct {
	agent      *loader.Agent
	minionName string
	username   string
	minion     loader.Minion
	client     *openrouter.Client
	messages   []openrouter.Message
}

func NewSession(agent *loader.Agent, minionName string, minion loader.Minion, client *openrouter.Client, username string) *Session {
	prompt := BuildPrompt(agent.Instinct)
	return &Session{
		agent:      agent,
		minionName: minionName,
		username:   username,
		minion:     minion,
		client:     client,
		messages: []openrouter.Message{
			{Role: "system", Content: prompt},
		},
	}
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
			fmt.Fprintf(&b, "## %s-%s\n\n%s\n\n", s.agent.Name, s.minionName, msg.Content)
		}
	}
	return os.WriteFile(filepath.Join(convsDir, name+".md"), []byte(b.String()), 0o644)
}

func (s *Session) Send(ctx context.Context, content string, onToken func(string)) (string, error) {
	s.messages = append(s.messages, openrouter.Message{
		Role:    "user",
		Content: content,
	})
	req := openrouter.NewRequest(s.minion, s.messages)
	stream, err := s.client.Complete(ctx, req)
	if err != nil {
		s.messages = s.messages[:len(s.messages)-1]
		return "", fmt.Errorf("completing request: %w", err)
	}
	defer stream.Close()
	var resp strings.Builder
	for {
		token, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.messages = s.messages[:len(s.messages)-1]
			return resp.String(), fmt.Errorf("streaming response: %w", err)
		}
		resp.WriteString(token)
		if onToken != nil {
			onToken(token)
		}
	}
	full := resp.String()
	s.messages = append(s.messages, openrouter.Message{
		Role:    "assistant",
		Content: full,
	})
	return full, nil
}
