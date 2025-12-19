package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultOpenAIModel = "gpt-4o-mini"

	envOpenAIAPIKey  = "OPENAI_API_KEY"
	envOpenAIModel   = "OPENAI_MODEL"
	envOpenAIBaseURL = "OPENAI_BASE_URL"
)

type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

func LoadOpenAIConfigFromEnv() (*OpenAIConfig, error) {
	key := strings.TrimSpace(os.Getenv(envOpenAIAPIKey))
	if key == "" {
		return nil, fmt.Errorf("config: %s is not set", envOpenAIAPIKey)
	}

	model := strings.TrimSpace(os.Getenv(envOpenAIModel))
	if model == "" {
		model = DefaultOpenAIModel
	}

	baseURL := strings.TrimSpace(os.Getenv(envOpenAIBaseURL))

	return &OpenAIConfig{
		APIKey:  key,
		Model:   model,
		BaseURL: baseURL,
	}, nil
}
