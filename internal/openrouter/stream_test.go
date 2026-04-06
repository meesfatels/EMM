package openrouter

import (
	"io"
	"strings"
	"testing"
)

func TestStreamRecv_ContentAndDone(t *testing.T) {
	body := strings.Join([]string{
		"data: {\"choices\":[{\"delta\":{\"content\":\"hello\"}}]}",
		"data:[DONE]",
		"",
	}, "\n")

	s := newStream(io.NopCloser(strings.NewReader(body)))
	defer s.Close()

	tok, err := s.Recv()
	if err != nil {
		t.Fatalf("Recv token error: %v", err)
	}
	if tok != "hello" {
		t.Fatalf("token = %q, want hello", tok)
	}

	_, err = s.Recv()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestStreamRecv_StreamErrorPayload(t *testing.T) {
	body := "data: {\"error\":{\"message\":\"rate limited\"}}\n"
	s := newStream(io.NopCloser(strings.NewReader(body)))
	defer s.Close()

	_, err := s.Recv()
	if err == nil || !strings.Contains(err.Error(), "rate limited") {
		t.Fatalf("expected stream error with message, got %v", err)
	}
}
