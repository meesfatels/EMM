package openrouter

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type Stream struct {
	reader  io.ReadCloser
	scanner *bufio.Scanner
}

func newStream(r io.ReadCloser) *Stream {
	return &Stream{
		reader:  r,
		scanner: bufio.NewScanner(r),
	}
}
func (s *Stream) Recv() (string, error) {
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return "", io.EOF
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
		return "", err
	}
	return "", io.EOF
}
func (s *Stream) Close() error {
	return s.reader.Close()
}
