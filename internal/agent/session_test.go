package agent_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/meesfatels/emm/internal/agent"
	"github.com/meesfatels/emm/internal/minion"
	"github.com/meesfatels/emm/internal/openrouter"
)

// sseServer returns a test server that streams the given text tokens as SSE.
func sseServer(tokens ...string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for _, tok := range tokens {
			fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"content\":%q}}]}\n\n", tok)
		}
		fmt.Fprintln(w, "data: [DONE]")
	}))
}

func newSession(t *testing.T, serverURL string) *agent.Session {
	t.Helper()
	a := &agent.Agent{Name: "test"}
	m := minion.Minion{"model": "test-model"}
	client := openrouter.NewClient("key", serverURL)
	return agent.NewSession(a, "test-minion", m, client, "user")
}

// ---- Send -------------------------------------------------------------------

func TestSend(t *testing.T) {
	t.Run("assembles streamed tokens", func(t *testing.T) {
		srv := sseServer("hello", " ", "world")
		defer srv.Close()

		result := newSession(t, srv.URL).Send(context.Background(), "hi", nil, nil)
		if result != "hello world" {
			t.Errorf("got %q, want %q", result, "hello world")
		}
	})

	t.Run("calls onToken for each chunk", func(t *testing.T) {
		srv := sseServer("a", "b", "c")
		defer srv.Close()

		var got []string
		newSession(t, srv.URL).Send(context.Background(), "hi",
			func(tok string) { got = append(got, tok) }, nil)
		if strings.Join(got, "") != "abc" {
			t.Errorf("onToken sequence: got %v", got)
		}
	})

	t.Run("appends messages to history", func(t *testing.T) {
		srv := sseServer("response")
		defer srv.Close()

		s := newSession(t, srv.URL)
		s.Send(context.Background(), "question", nil, nil)

		msgs := s.Messages()
		// system + user + assistant
		if len(msgs) != 3 {
			t.Fatalf("expected 3 messages, got %d", len(msgs))
		}
		if msgs[1].Role != "user" || msgs[1].Content != "question" {
			t.Errorf("user message: %+v", msgs[1])
		}
		if msgs[2].Role != "assistant" || msgs[2].Content != "response" {
			t.Errorf("assistant message: %+v", msgs[2])
		}
	})
}

// ---- Save / Load ------------------------------------------------------------

func TestSaveLoad(t *testing.T) {
	srv := sseServer("hello")
	defer srv.Close()

	dir := t.TempDir()
	s := newSession(t, srv.URL)
	s.Send(context.Background(), "ping", nil, nil)
	s.Save(dir, "myconv")

	s2 := newSession(t, srv.URL)
	if !s2.Load(dir, "myconv") {
		t.Fatal("Load returned false")
	}

	orig := s.Messages()
	loaded := s2.Messages()
	if len(orig) != len(loaded) {
		t.Fatalf("message count: orig %d, loaded %d", len(orig), len(loaded))
	}
	for i := range orig {
		if orig[i].Role != loaded[i].Role || orig[i].Content != loaded[i].Content {
			t.Errorf("message[%d]: got {%s %q}, want {%s %q}",
				i, loaded[i].Role, loaded[i].Content, orig[i].Role, orig[i].Content)
		}
	}
}

func TestLoad_NotFound(t *testing.T) {
	s := newSession(t, "http://localhost")
	if s.Load(t.TempDir(), "nonexistent") {
		t.Error("Load should return false for missing conversation")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	srv := sseServer("hi")
	defer srv.Close()

	dir := t.TempDir()
	s := newSession(t, srv.URL)
	s.Send(context.Background(), "test", nil, nil)
	s.Save(dir, "saved")

	data, err := os.ReadFile(dir + "/conversations/saved.md")
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if !strings.Contains(string(data), "agent: test") {
		t.Errorf("saved file missing agent header:\n%s", data)
	}
}
