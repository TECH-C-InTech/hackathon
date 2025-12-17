package repository

import (
	"context"

	"backend/internal/domain/post"
)

// PostRepository は闇投稿を永続化するための抽象.
type PostRepository interface {
	// FindPending は pending 状態の投稿を 1 件返す.
	// 見つからなかった場合は (nil, nil) を返す。
	FindPending(ctx context.Context) (*post.Post, error)
	// Update は投稿の内容や状態を保存する。
	Update(ctx context.Context, p *post.Post) error
}
