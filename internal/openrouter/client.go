package openrouter

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Complete(ctx context.Context, req *Request) (*Stream, error) {
	body, err := req.encode()
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openrouter %d: %s", resp.StatusCode, string(errBody))
	}
	return newStream(resp.Body), nil
}
