package queue

import "context"

// JobQueue は Cloud Tasks などのジョブキューからワーカーを起動するための抽象.
type JobQueue interface {
	RunWorker(ctx context.Context, handler Handler) error
}

// Handler は受け取ったジョブを処理する関数.
type Handler interface {
	Handle(ctx context.Context, payload []byte) error
}

// HandlerFunc は Handler を関数で表現するためのアダプタ.
type HandlerFunc func(ctx context.Context, payload []byte) error

// Handle は f(ctx, payload) を呼び出す。
func (f HandlerFunc) Handle(ctx context.Context, payload []byte) error {
	return f(ctx, payload)
}
