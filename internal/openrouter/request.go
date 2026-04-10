package openrouter

import (
	"bytes"
	"encoding/json"
	"io"
)

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Tool struct {
	Type     string         `json:"type"`
	Function ToolDefinition `json:"function"`
}

type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type Request struct {
	Minion   map[string]any
	Messages []Message
	Tools    []Tool
}

func (r *Request) encode() io.Reader {
	payload := make(map[string]any, len(r.Minion)+3)
	for k, v := range r.Minion {
		payload[k] = v
	}
	payload["messages"] = r.Messages
	payload["stream"] = true
	if len(r.Tools) > 0 {
		payload["tools"] = r.Tools
	}
	data, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}
