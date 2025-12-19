package config

import (
	"os"
	"strings"
)

const envLLMProvider = "LLM_PROVIDER"

/**
 * LLM_PROVIDER 環境変数から使用する LLM 名を取得し、未設定時は openai を返す。
 */
func LoadLLMProvider() string {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv(envLLMProvider)))
	if provider == "" {
		return "openai"
	}
	switch provider {
	case "openai", "gemini":
		return provider
	default:
		return "openai"
	}
}
