package draw

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	drawdomain "backend/internal/domain/draw"
	"backend/internal/domain/post"
	"backend/internal/port/repository"
)

func TestDrawFortune_Success(t *testing.T) {
	t.Parallel()

	repo := &fakeDrawRepository{
		draws: []*drawdomain.Draw{
			newVerifiedDraw(t, "post-1", "fortune-1"),
			newVerifiedDraw(t, "post-2", "fortune-2"),
		},
	}
	usecase := NewFortuneUsecase(repo)
	usecase.rand = rand.New(rand.NewSource(1)) // テストの乱択結果を固定

	got, err := usecase.DrawFortune(context.Background())
	if err != nil {
		t.Fatalf("DrawFortune() error = %v", err)
	}
	if got == nil {
		t.Fatal("DrawFortune() returned nil draw")
	}
	if got.Status() != drawdomain.StatusVerified {
		t.Fatalf("expected verified status, got %s", got.Status())
	}
	if got.PostID() != post.DarkPostID("post-2") {
		t.Fatalf("unexpected draw selected, got %s", got.PostID())
	}
}

func TestDrawFortune_EmptyResults(t *testing.T) {
	t.Parallel()

	usecase := NewFortuneUsecase(&fakeDrawRepository{})
	_, err := usecase.DrawFortune(context.Background())
	if !errors.Is(err, drawdomain.ErrEmptyResult) {
		t.Fatalf("expected ErrEmptyResult, got %v", err)
	}
}

func TestDrawFortune_RepositoryError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("repository failure")
	usecase := NewFortuneUsecase(&fakeDrawRepository{listErr: expectedErr})

	_, err := usecase.DrawFortune(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

type fakeDrawRepository struct {
	repository.DrawRepository

	draws   []*drawdomain.Draw
	listErr error
}

func (f *fakeDrawRepository) Create(ctx context.Context, d *drawdomain.Draw) error {
	return nil
}

func (f *fakeDrawRepository) GetByPostID(ctx context.Context, postID post.DarkPostID) (*drawdomain.Draw, error) {
	return nil, repository.ErrDrawNotFound
}

func (f *fakeDrawRepository) ListReady(ctx context.Context) ([]*drawdomain.Draw, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.draws, nil
}

func newVerifiedDraw(t *testing.T, postID, result string) *drawdomain.Draw {
	t.Helper()

	d, err := drawdomain.New(post.DarkPostID(postID), drawdomain.FormattedContent(result))
	if err != nil {
		t.Fatalf("drawdomain.New() error = %v", err)
	}
	d.MarkVerified()
	return d
}
