package openrouter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Stream struct {
	reader  io.ReadCloser
	scanner *bufio.Scanner
}

func newStream(r io.ReadCloser) *Stream {
	scanner := bufio.NewScanner(r)
	// SSE chunks can exceed the scanner default (64K) for large deltas.
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	return &Stream{
		reader:  r,
		scanner: scanner,
	}
}
func (s *Stream) Recv() (string, error) {
	for s.scanner.Scan() {
		line := strings.TrimSpace(s.scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}
		if data == "[DONE]" {
			return "", io.EOF
		}

		var errChunk struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(data), &errChunk); err == nil && errChunk.Error.Message != "" {
			return "", fmt.Errorf("openrouter stream error: %s", errChunk.Error.Message)
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			return chunk.Choices[0].Delta.Content, nil
		}
	}
	if err := s.scanner.Err(); err != nil {
		return "", fmt.Errorf("reading stream: %w", err)
	}
	return "", io.EOF
}
func (s *Stream) Close() error {
	return s.reader.Close()
}
