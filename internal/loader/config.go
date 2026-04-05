package loader
import (
	"fmt"
)
type Config map[string]any
func (c Config) APIKey() (string, error) {
	key, ok := c["api_key"]
	if !ok {
		return "", fmt.Errorf("api_key not set in emm.yaml")
	}
	s, ok := key.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("api_key is empty in emm.yaml")
	}
	return s, nil
}
func (c Config) BaseURL() string {
	url, ok := c["base_url"]
	if !ok {
		return "https://openrouter.ai/api/v1/chat/completions"
	}
	s, ok := url.(string)
	if !ok || s == "" {
		return "https://openrouter.ai/api/v1/chat/completions"
	}
	return s
}
