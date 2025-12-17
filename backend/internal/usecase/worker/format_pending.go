package worker

import (
	"context"
	"fmt"

	"backend/internal/port/llm"
	"backend/internal/port/repository"
)

// FormatPending は pending 投稿を LLM に投げて ready へ遷移させるユースケース。
type FormatPending struct {
	postRepo  repository.PostRepository
	formatter llm.Formatter
}

// NewFormatPending は FormatPending を生成する。
func NewFormatPending(postRepo repository.PostRepository, formatter llm.Formatter) *FormatPending {
	return &FormatPending{
		postRepo:  postRepo,
		formatter: formatter,
	}
}

// Run は pending 投稿を 1 件処理する。
//
// リトライ方針:
// - Repository / Formatter の呼び出しが失敗した場合はエラーを返し、呼び出し元（Cloud Tasks 等）の再実行に任せる
// - pending 投稿が見つからない場合は何もせずに終了し、次のジョブ／ポーリングで再度実行する
func (uc *FormatPending) Run(ctx context.Context) error {
	p, err := uc.postRepo.FindPending(ctx)
	if err != nil {
		return fmt.Errorf("find pending post: %w", err)
	}
	if p == nil {
		return nil
	}

	formatted, err := uc.formatter.Format(ctx, p.Content())
	if err != nil {
		return fmt.Errorf("format content: %w", err)
	}
	if err := p.UpdateContent(formatted); err != nil {
		return fmt.Errorf("update content: %w", err)
	}
	if err := p.MarkReady(); err != nil {
		return fmt.Errorf("mark ready: %w", err)
	}
	if err := uc.postRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("update post: %w", err)
	}

	return nil
}
