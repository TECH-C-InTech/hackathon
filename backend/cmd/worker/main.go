package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/adapter/llm/gemini"
	"backend/internal/adapter/repository/memory"
	"backend/internal/app"
	"backend/internal/domain/post"
	workerusecase "backend/internal/usecase/worker"
)

/**
 * 非同期処理する別プロセスのエントリポイント
 */
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo := memory.NewPostRepository()
	seedFromEnv(repo)

	formatter := gemini.NewFormatter()
	formatPending := workerusecase.NewFormatPending(repo, formatter)

	workerApp := app.NewWorkerApp(formatPending, nil, log.Default())

	mode := os.Getenv("WORKER_MODE")
	if mode == "" {
		mode = string(app.WorkerModePolling)
	}

	if err := workerApp.Run(ctx, app.WorkerMode(mode)); err != nil {
		log.Fatalf("worker stopped: %v", err)
	}
}

func seedFromEnv(repo *memory.PostRepository) {
	content := os.Getenv("SEED_POST_CONTENT")
	if content == "" {
		return
	}

	id := fmt.Sprintf("seed-%d", time.Now().UnixNano())
	p, err := post.New(id, content)
	if err != nil {
		log.Printf("worker: failed to seed pending post: %v", err)
		return
	}

	repo.Insert(p)
	log.Printf("worker: seeded pending post %s", id)
}
