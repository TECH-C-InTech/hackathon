package post

import (
	"context"
	"errors"

	"backend/internal/domain/post"
	"backend/internal/port/queue"
	"backend/internal/port/repository"
)

var (
	ErrNilInput            = errors.New("create_post: 入力が指定されていません")
	ErrPostAlreadyExists   = errors.New("create_post: 投稿がすでに存在します")
	ErrJobAlreadyScheduled = errors.New("create_post: 整形ジョブがすでに登録済みです")
)

// 闇投稿作成の入力値
type CreatePostInput struct {
	DarkPostID string
	Content    string
}

// 闇投稿作成後に呼び出し側へ返す値
type CreatePostOutput struct {
	DarkPostID string
}

/**
 * 闇投稿作成のユースケース
 * postRepo: 投稿リポジトリ
 * jobQueue: 整形ジョブキュー
 */
type CreatePostUsecase struct {
	postRepo repository.PostRepository
	jobQueue queue.JobQueue
}

/**
 * ユースケース毎に初期化
 */
func NewCreatePostUsecase(postRepo repository.PostRepository, jobQueue queue.JobQueue) *CreatePostUsecase {
	return &CreatePostUsecase{
		postRepo: postRepo,
		jobQueue: jobQueue,
	}
}

/**
 * 闇投稿作成の実行
 */
func (u *CreatePostUsecase) Execute(ctx context.Context, in *CreatePostInput) (*CreatePostOutput, error) {
	// 空かどうか
	if in == nil {
		return nil, ErrNilInput
	}

	// 投稿オブジェクトの生成
	p, err := post.New(post.DarkPostID(in.DarkPostID), post.DarkContent(in.Content))
	if err != nil {
		return nil, err
	}

	// 投稿の保存
	if err := u.postRepo.Create(ctx, p); err != nil {
		// 重複時はエラー
		if errors.Is(err, repository.ErrPostAlreadyExists) {
			return nil, ErrPostAlreadyExists
		}

		return nil, err
	}

	// 整形ジョブの登録
	if err := u.jobQueue.EnqueueFormat(ctx, p.ID()); err != nil {
		// 重複時はエラー
		if errors.Is(err, queue.ErrJobAlreadyScheduled) {
			return nil, ErrJobAlreadyScheduled
		}

		return nil, err
	}

	return &CreatePostOutput{DarkPostID: string(p.ID())}, nil
}
