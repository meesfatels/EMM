package openrouter

import (
	"context"
	"io"
	"net/http"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{apiKey: apiKey, baseURL: baseURL, httpClient: &http.Client{}}
}

func (c *Client) Complete(ctx context.Context, req *Request) *Stream {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, req.encode())
	if err != nil {
		panic(err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		panic("openrouter: " + string(body))
	}
	return NewStream(resp.Body)
}
