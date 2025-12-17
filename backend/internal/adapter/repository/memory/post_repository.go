package memory

import (
	"context"
	"sync"

	"backend/internal/domain/post"
	"backend/internal/port/repository"
)

var _ repository.PostRepository = (*PostRepository)(nil)

// PostRepository はローカル開発向けの簡易インメモリ実装。
type PostRepository struct {
	mu    sync.Mutex
	posts []*post.Post
}

// NewPostRepository は初期化する。
func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: make([]*post.Post, 0),
	}
}

// Insert は pending 投稿を追加する（デバッグ用）。
func (r *PostRepository) Insert(p *post.Post) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts = append(r.posts, clonePost(p))
}

// FindPending は pending な投稿を取得する。
func (r *PostRepository) FindPending(context.Context) (*post.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.posts {
		if p.Status() == post.StatusPending {
			return clonePost(p), nil
		}
	}
	return nil, nil
}

// Update は投稿を保存する。
func (r *PostRepository) Update(ctx context.Context, updated *post.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.posts {
		if p.ID() == updated.ID() {
			r.posts[i] = clonePost(updated)
			return nil
		}
	}

	r.posts = append(r.posts, clonePost(updated))
	return nil
}

func clonePost(p *post.Post) *post.Post {
	cloned, err := post.Restore(p.ID(), p.Content(), p.Status())
	if err != nil {
		// Restore only fails for invalid data, which should never happen here.
		panic(err)
	}
	return cloned
}
