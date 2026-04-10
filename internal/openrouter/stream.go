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
	assembled map[int]*ToolCall
	toolCalls []ToolCall
}

func NewStream(r io.ReadCloser) *Stream {
	s := &Stream{
		reader:    r,
		scanner:   bufio.NewScanner(r),
		assembled: make(map[int]*ToolCall),
	}
	s.scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return s
}

// Recv returns the next text token and true, or ("", false) when the stream ends.
// Tool call fragments are collected internally; call ToolCalls() after false is returned.
func (s *Stream) Recv() (string, bool) {
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			s.toolCalls = s.finishToolCalls()
			return "", false
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
				a = &ToolCall{}
				s.assembled[tc.Index] = a
			}
			if tc.ID != "" {
				a.ID = tc.ID
			}
			if tc.Type != "" {
				a.Type = tc.Type
			}
			if tc.Function.Name != "" {
				a.Function.Name = tc.Function.Name
			}
			a.Function.Arguments += tc.Function.Arguments
		}

		if delta.Content != "" {
			return delta.Content, true
		}
	}
	s.toolCalls = s.finishToolCalls()
	return "", false
}

func (s *Stream) ToolCalls() []ToolCall {
	return s.toolCalls
}

func (s *Stream) Close() {
	s.reader.Close()
}

func (s *Stream) finishToolCalls() []ToolCall {
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
		result = append(result, *s.assembled[k])
	}
	return result
}
