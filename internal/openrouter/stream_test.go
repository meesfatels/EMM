package openrouter_test

import (
	"io"
	"strings"
	"testing"

	"github.com/meesfatels/emm/internal/openrouter"
)

// stream builds a Stream from raw SSE lines (no blank-line separators needed).
func stream(lines ...string) *openrouter.Stream {
	body := strings.Join(lines, "\n") + "\n"
	return openrouter.NewStream(io.NopCloser(strings.NewReader(body)))
}

// collect drains a stream and returns all tokens.
func collect(s *openrouter.Stream) []string {
	var tokens []string
	for {
		tok, ok := s.Recv()
		if !ok {
			break
		}
		if tok != "" {
			tokens = append(tokens, tok)
		}
	}
	return tokens
}

func TestRecv_Text(t *testing.T) {
	s := stream(
		`data: {"choices":[{"delta":{"content":"hello"}}]}`,
		`data: {"choices":[{"delta":{"content":" world"}}]}`,
		`data: [DONE]`,
	)
	got := strings.Join(collect(s), "")
	if got != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestRecv_Done(t *testing.T) {
	s := stream(`data: [DONE]`)
	tokens := collect(s)
	if len(tokens) != 0 {
		t.Errorf("expected no tokens after [DONE], got %v", tokens)
	}
}

func TestRecv_SkipsNonDataLines(t *testing.T) {
	s := stream(
		`event: message`,
		`id: 1`,
		`data: {"choices":[{"delta":{"content":"ok"}}]}`,
		`data: [DONE]`,
	)
	got := strings.Join(collect(s), "")
	if got != "ok" {
		t.Errorf("got %q, want %q", got, "ok")
	}
}

func TestRecv_SkipsInvalidJSON(t *testing.T) {
	s := stream(
		`data: not json at all`,
		`data: {"choices":[{"delta":{"content":"fine"}}]}`,
		`data: [DONE]`,
	)
	got := strings.Join(collect(s), "")
	if got != "fine" {
		t.Errorf("got %q, want %q", got, "fine")
	}
}

func TestRecv_ToolCalls(t *testing.T) {
	// Tool call spread across multiple chunks, as the real API sends it.
	s := stream(
		`data: {"choices":[{"delta":{"tool_calls":[{"index":0,"id":"tc1","type":"function","function":{"name":"run_shell","arguments":""}}]}}]}`,
		`data: {"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"cmd\":"}}]}}]}`,
		`data: {"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"ls\"}"}}]}}]}`,
		`data: [DONE]`,
	)

	tokens := collect(s)
	if len(tokens) != 0 {
		t.Errorf("expected no text tokens for tool-call response, got %v", tokens)
	}

	calls := s.ToolCalls()
	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}
	tc := calls[0]
	if tc.ID != "tc1" {
		t.Errorf("ID: got %q, want %q", tc.ID, "tc1")
	}
	if tc.Function.Name != "run_shell" {
		t.Errorf("Name: got %q, want %q", tc.Function.Name, "run_shell")
	}
	if tc.Function.Arguments != `{"cmd":"ls"}` {
		t.Errorf("Arguments: got %q, want %q", tc.Function.Arguments, `{"cmd":"ls"}`)
	}
}

func TestRecv_MultipleToolCalls(t *testing.T) {
	s := stream(
		`data: {"choices":[{"delta":{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"tool_a","arguments":"{}"}}]}}]}`,
		`data: {"choices":[{"delta":{"tool_calls":[{"index":1,"id":"b","type":"function","function":{"name":"tool_b","arguments":"{}"}}]}}]}`,
		`data: [DONE]`,
	)
	collect(s)

	calls := s.ToolCalls()
	if len(calls) != 2 {
		t.Fatalf("expected 2 tool calls, got %d", len(calls))
	}
	// Must be ordered by index.
	if calls[0].Function.Name != "tool_a" || calls[1].Function.Name != "tool_b" {
		t.Errorf("wrong order or names: %v", calls)
	}
}
