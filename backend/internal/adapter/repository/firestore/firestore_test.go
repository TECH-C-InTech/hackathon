package firestore

import (
	"context"
	"os"
	"testing"

	drawdomain "backend/internal/domain/draw"
	"backend/internal/domain/post"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const (
	testProjectID = "firestore-integration-test"
)

// newTestFirestoreClient は Firestore エミュレータに接続するクライアントを返す。
// 環境変数が未設定の場合はテストをスキップする。
func newTestFirestoreClient(t *testing.T) *firestore.Client {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set; skipping Firestore integration tests")
	}
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = testProjectID
	}
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})
	return client
}

// truncateCollection は指定コレクションの全ドキュメントを削除する。
func truncateCollection(t *testing.T, client *firestore.Client, collection string) {
	t.Helper()
	ctx := context.Background()
	iter := client.Collection(collection).Documents(ctx)
	batch := client.Batch()
	count := 0
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			t.Fatalf("iterate %s: %v", collection, err)
		}
		batch.Delete(doc.Ref)
		count++
		if count == 500 {
			if _, err := batch.Commit(ctx); err != nil {
				t.Fatalf("commit delete batch: %v", err)
			}
			batch = client.Batch()
			count = 0
		}
	}
	if count > 0 {
		if _, err := batch.Commit(ctx); err != nil {
			t.Fatalf("commit delete batch: %v", err)
		}
	}
}
func TestDrawRepository_Integration(t *testing.T) {
	client := newTestFirestoreClient(t)
	truncateCollection(t, client, drawsCollection)

	repo, err := NewDrawRepository(client)
	if err != nil {
		t.Fatalf("new draw repo: %v", err)
	}

	ctx := context.Background()
	draw, err := drawdomain.New(post.DarkPostID("post-1"), drawdomain.FormattedContent("fortune smiles"))
	if err != nil {
		t.Fatalf("new draw: %v", err)
	}
	draw.MarkVerified()

	if err := repo.Create(ctx, draw); err != nil {
		t.Fatalf("create draw: %v", err)
	}

	fetched, err := repo.GetByPostID(ctx, "post-1")
	if err != nil {
		t.Fatalf("get draw: %v", err)
	}
	if fetched.Result() != draw.Result() || fetched.Status() != draw.Status() {
		t.Fatalf("fetched draw mismatch")
	}

	list, err := repo.ListReady(ctx)
	if err != nil {
		t.Fatalf("list ready: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 draw got %d", len(list))
	}
}
