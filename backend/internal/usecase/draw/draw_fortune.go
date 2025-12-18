package draw

import (
	"context"
	"math/rand"
	"time"

	drawdomain "backend/internal/domain/draw"
	"backend/internal/port/repository"
)

// FortuneUsecase は検証済みのおみくじを 1 件返すユースケース。
type FortuneUsecase struct {
	repo repository.DrawRepository
	rand *rand.Rand
}

// NewFortuneUsecase は FortuneUsecase を生成する。
func NewFortuneUsecase(repo repository.DrawRepository) *FortuneUsecase {
	return &FortuneUsecase{
		repo: repo,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// DrawFortune は Verified 状態のおみくじから 1 件をランダムに返す。
func (u *FortuneUsecase) DrawFortune(ctx context.Context) (*drawdomain.Draw, error) {
	draws, err := u.repo.ListReady(ctx)
	if err != nil {
		return nil, err
	}

	verified := make([]*drawdomain.Draw, 0, len(draws))
	for _, d := range draws {
		if d == nil {
			continue
		}
		if d.Status() != drawdomain.StatusVerified {
			continue
		}
		verified = append(verified, d)
	}

	if len(verified) == 0 {
		return nil, drawdomain.ErrEmptyResult
	}
	if len(verified) == 1 {
		return verified[0], nil
	}

	index := u.rand.Intn(len(verified))
	return verified[index], nil
}
