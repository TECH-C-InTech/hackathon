package llm

import "context"

// Formatter は LLM で投稿本文を整形するためのインターフェース。
type Formatter interface {
	Format(ctx context.Context, content string) (string, error)
}
