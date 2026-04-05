package openrouter

import (
	"bytes"
	"encoding/json"
	"io"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Request struct {
	Minion   map[string]any
	Messages []Message
}

func NewRequest(minion map[string]any, messages []Message) *Request {
	return &Request{
		Minion:   minion,
		Messages: messages,
	}
}
func (r *Request) encode() (io.Reader, error) {
	payload := make(map[string]any, len(r.Minion)+2)
	for k, v := range r.Minion {
		payload[k] = v
	}
	payload["messages"] = r.Messages
	payload["stream"] = true
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
