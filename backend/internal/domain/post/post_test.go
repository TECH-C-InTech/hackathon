package post

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()

	post, err := New("post-id", "闇がおおい")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.ID() != "post-id" {
		t.Fatalf("unexpected id: %s", post.ID())
	}
	if post.Content() != "闇がおおい" {
		t.Fatalf("unexpected content: %s", post.Content())
	}
	if post.Status() != StatusPending {
		t.Fatalf("expected pending but got %s", post.Status())
	}
}

func TestNew_EmptyContent(t *testing.T) {
	t.Parallel()

	if _, err := New("post-id", ""); err != ErrEmptyContent {
		t.Fatalf("expected ErrEmptyContent but got %v", err)
	}
}

func TestMarkReady(t *testing.T) {
	t.Parallel()

	post, err := New("id", "闇")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := post.MarkReady(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.Status() != StatusReady {
		t.Fatalf("expected ready but got %s", post.Status())
	}
}

func TestMarkReady_InvalidTransition(t *testing.T) {
	t.Parallel()

	post, err := Restore("id", "闇", StatusReady)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := post.MarkReady(); err != ErrInvalidStatusTransition {
		t.Fatalf("expected ErrInvalidStatusTransition but got %v", err)
	}
}

func TestRestore_InvalidStatus(t *testing.T) {
	t.Parallel()

	if _, err := Restore("id", "闇", Status("unknown")); err != ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus but got %v", err)
	}
}

func TestUpdateContent(t *testing.T) {
	t.Parallel()

	p, err := New("id", "before")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.UpdateContent("after"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Content() != "after" {
		t.Fatalf("expected updated content but got %s", p.Content())
	}
}

func TestUpdateContent_Empty(t *testing.T) {
	t.Parallel()

	p, err := New("id", "闇")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.UpdateContent(""); err != ErrEmptyContent {
		t.Fatalf("expected ErrEmptyContent but got %v", err)
	}
}
