package app

import (
	"context"
	"errors"
	"log"
	"time"

	"backend/internal/port/queue"
	usecaseworker "backend/internal/usecase/worker"
)

// WorkerMode はワーカーの実行モード。
type WorkerMode string

const (
	// WorkerModeQueue はジョブキュー経由の実行。
	WorkerModeQueue WorkerMode = "queue"
	// WorkerModePolling は DB を定期ポーリングして自己完結的に実行。
	WorkerModePolling WorkerMode = "polling"

	defaultPollInterval = 5 * time.Second
)

// WorkerApp は LLM 整形ワーカーのエントリポイント。
type WorkerApp struct {
	formatPending *usecaseworker.FormatPending
	queue         queue.JobQueue
	logger        *log.Logger
	pollInterval  time.Duration
}

// WorkerOption は WorkerApp のオプション。
type WorkerOption func(*WorkerApp)

// WithPollInterval はポーリング間隔を設定する。
func WithPollInterval(d time.Duration) WorkerOption {
	return func(app *WorkerApp) {
		if d > 0 {
			app.pollInterval = d
		}
	}
}

// NewWorkerApp は WorkerApp を生成する。
func NewWorkerApp(formatPending *usecaseworker.FormatPending, q queue.JobQueue, logger *log.Logger, opts ...WorkerOption) *WorkerApp {
	if logger == nil {
		logger = log.Default()
	}
	app := &WorkerApp{
		formatPending: formatPending,
		queue:         q,
		logger:        logger,
		pollInterval:  defaultPollInterval,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// Run は指定モードでワーカーを起動する。
func (a *WorkerApp) Run(ctx context.Context, mode WorkerMode) error {
	switch mode {
	case WorkerModeQueue:
		if a.queue == nil {
			return errors.New("worker queue not configured")
		}
		return a.queue.RunWorker(ctx, queue.HandlerFunc(func(ctx context.Context, _ []byte) error {
			return a.runOnce(ctx)
		}))
	case WorkerModePolling:
		return a.runPolling(ctx)
	default:
		return errors.New("unknown worker mode")
	}
}

func (a *WorkerApp) runPolling(ctx context.Context) error {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	for {
		_ = a.runOnce(ctx)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (a *WorkerApp) runOnce(ctx context.Context) error {
	err := a.formatPending.Run(ctx)
	if err != nil {
		a.logger.Printf("worker: failed to process pending post: %v", err)
	}
	return err
}
