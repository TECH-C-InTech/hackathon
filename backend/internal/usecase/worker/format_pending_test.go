package worker

import (
	"context"
	"errors"
	"testing"

	"backend/internal/domain/post"
)

func TestFormatPending_RunSuccess(t *testing.T) {
	t.Parallel()

	p, err := post.New("id", "raw")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo := &fakePostRepository{
		pending: p,
	}
	formatter := &fakeFormatter{
		output: "formatted",
	}

	usecase := NewFormatPending(repo, formatter)

	if err := usecase.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.updated {
		t.Fatal("expected update to be called")
	}
	if repo.savedPost.Content() != "formatted" {
		t.Fatalf("expected formatted content but got %s", repo.savedPost.Content())
	}
	if repo.savedPost.Status() != post.StatusReady {
		t.Fatalf("expected status ready but got %s", repo.savedPost.Status())
	}
}

func TestFormatPending_NoPending(t *testing.T) {
	t.Parallel()

	repo := &fakePostRepository{}
	usecase := NewFormatPending(repo, &fakeFormatter{})

	if err := usecase.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updated {
		t.Fatal("update should not be called when no pending post")
	}
}

func TestFormatPending_FormatError(t *testing.T) {
	t.Parallel()

	p, err := post.New("id", "raw")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	usecase := NewFormatPending(&fakePostRepository{pending: p}, &fakeFormatter{err: errors.New("boom")})

	if err := usecase.Run(context.Background()); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestFormatPending_UpdateError(t *testing.T) {
	t.Parallel()

	p, err := post.New("id", "raw")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo := &fakePostRepository{
		pending:   p,
		updateErr: errors.New("db down"),
	}

	usecase := NewFormatPending(repo, &fakeFormatter{})

	if err := usecase.Run(context.Background()); err == nil {
		t.Fatal("expected error but got nil")
	}
}

type fakePostRepository struct {
	pending   *post.Post
	findErr   error
	updateErr error

	updated   bool
	savedPost *post.Post
}

func (r *fakePostRepository) FindPending(context.Context) (*post.Post, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.pending, nil
}

func (r *fakePostRepository) Update(ctx context.Context, p *post.Post) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	r.updated = true
	r.savedPost = p
	return nil
}

type fakeFormatter struct {
	output string
	err    error
}

func (f *fakeFormatter) Format(context.Context, string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if f.output != "" {
		return f.output, nil
	}
	return "formatted", nil
}
