package handler

import (
	"context"
	"errors"
	"net/http"

	drawdomain "backend/internal/domain/draw"

	"github.com/gin-gonic/gin"
)

const (
	messageDrawsEmpty    = "no verified draws available"
	messageInternalError = "internal server error"
)

// FortuneUsecase は検証済みのおみくじを 1 件返すユースケースの契約。
type FortuneUsecase interface {
	DrawFortune(ctx context.Context) (*drawdomain.Draw, error)
}

// DrawHandler はおみくじ関連の HTTP ハンドラをまとめる。
type DrawHandler struct {
	usecase FortuneUsecase
}

// NewDrawHandler は DrawHandler を生成する。
func NewDrawHandler(usecase FortuneUsecase) *DrawHandler {
	return &DrawHandler{usecase: usecase}
}

// DrawResponse は GET /draws/random のレスポンス。
type DrawResponse struct {
	PostID string `json:"post_id"`
	Result string `json:"result"`
	Status string `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
}

// GetRandomDraw は Verified な結果を 1 件ランダムに返す。
func (h *DrawHandler) GetRandomDraw(c *gin.Context) {
	draw, err := h.usecase.DrawFortune(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, DrawResponse{
		PostID: string(draw.PostID()),
		Result: string(draw.Result()),
		Status: string(draw.Status()),
	})
}

func (h *DrawHandler) handleError(c *gin.Context, err error) {
	if errors.Is(err, drawdomain.ErrEmptyResult) {
		c.JSON(http.StatusNotFound, errorResponse{Message: messageDrawsEmpty})
		return
	}
	c.JSON(http.StatusInternalServerError, errorResponse{Message: messageInternalError})
}
