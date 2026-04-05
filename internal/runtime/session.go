package runtime

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/meesfatels/emm/internal/loader"
	"github.com/meesfatels/emm/internal/openrouter"
)

type Session struct {
	agent    *loader.Agent
	minion   loader.Minion
	client   *openrouter.Client
	messages []openrouter.Message
}

func NewSession(agent *loader.Agent, minion loader.Minion, client *openrouter.Client) *Session {
	prompt := BuildPrompt(agent.Instinct)
	return &Session{
		agent:  agent,
		minion: minion,
		client: client,
		messages: []openrouter.Message{
			{Role: "system", Content: prompt},
		},
	}
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
	var b strings.Builder
	for {
		token, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.messages = s.messages[:len(s.messages)-1]
			return b.String(), fmt.Errorf("streaming response: %w", err)
		}
		b.WriteString(token)
		if onToken != nil {
			onToken(token)
		}
	}
	full := b.String()
	s.messages = append(s.messages, openrouter.Message{
		Role:    "assistant",
		Content: full,
	})
	return full, nil
}
