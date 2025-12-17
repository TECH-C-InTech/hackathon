package gemini

import (
	"context"
	"log"

	"backend/internal/port/llm"
)

var _ llm.Formatter = (*Formatter)(nil)

// Formatter は Gemini 呼び出しのダミー実装。
// 実プロダクションではここで本物の Gemini API を呼び出す。
type Formatter struct{}

// NewFormatter は Formatter を生成する。
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Format は与えられたテキストを整形する。
// ひとまずパススルー実装としてログのみ出力する。
func (f *Formatter) Format(ctx context.Context, content string) (string, error) {
	log.Printf("gemini formatter called (len=%d)", len(content))
	return content, nil
}
