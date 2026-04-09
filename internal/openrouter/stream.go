package openrouter

import (
	"bufio"
	"encoding/json"
	"io"
	"slices"
	"strings"
)

type Stream struct {
	reader    io.ReadCloser
	scanner   *bufio.Scanner
	assembled map[int]*toolCallAssembler
	toolCalls []ToolCall
}

func newStream(r io.ReadCloser) *Stream {
	s := &Stream{
		reader:    r,
		scanner:   bufio.NewScanner(r),
		assembled: make(map[int]*toolCallAssembler),
	}
	s.scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return s
}

// Recv returns the next text token. Returns io.EOF when the stream ends.
// Tool call fragments are collected internally; call ToolCalls() after EOF.
func (s *Stream) Recv() (string, error) {
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			s.toolCalls = s.assembleToolCalls()
			return "", io.EOF
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content   string `json:"content"`
					ToolCalls []struct {
						Index    int    `json:"index"`
						ID       string `json:"id"`
						Type     string `json:"type"`
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil || len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta

		for _, tc := range delta.ToolCalls {
			a, ok := s.assembled[tc.Index]
			if !ok {
				a = &toolCallAssembler{}
				s.assembled[tc.Index] = a
			}
			if tc.ID != "" {
				a.id = tc.ID
			}
			if tc.Type != "" {
				a.typ = tc.Type
			}
			if tc.Function.Name != "" {
				a.name = tc.Function.Name
			}
			a.args += tc.Function.Arguments
		}

		if delta.Content != "" {
			return delta.Content, nil
		}
	}
	if err := s.scanner.Err(); err != nil {
		return "", err
	}
	s.toolCalls = s.assembleToolCalls()
	return "", io.EOF
}

// ToolCalls returns the tool calls collected during streaming.
// Only valid after Recv returns io.EOF.
func (s *Stream) ToolCalls() []ToolCall {
	return s.toolCalls
}

func (s *Stream) Close() error {
	return s.reader.Close()
}

type toolCallAssembler struct {
	id   string
	typ  string
	name string
	args string
}

func (s *Stream) assembleToolCalls() []ToolCall {
	if len(s.assembled) == 0 {
		return nil
	}
	keys := make([]int, 0, len(s.assembled))
	for k := range s.assembled {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	result := make([]ToolCall, 0, len(keys))
	for _, k := range keys {
		a := s.assembled[k]
		result = append(result, ToolCall{
			ID:   a.id,
			Type: a.typ,
			Function: ToolFunction{
				Name:      a.name,
				Arguments: a.args,
			},
		})
	}
	return result
}
