package handler

import (
	"errors"
	"net/http"
	"strings"

	postdomain "backend/internal/domain/post"
	postusecase "backend/internal/usecase/post"

	"github.com/gin-gonic/gin"
)

const (
	messagePostInvalidRequest = "invalid post request"
	messagePostConflict       = "post already exists"
)

// PostHandler は投稿関連の HTTP ハンドラをまとめる。
type PostHandler struct {
	createUsecase *postusecase.CreatePostUsecase
}

// NewPostHandler は PostHandler を生成する。
func NewPostHandler(usecase *postusecase.CreatePostUsecase) *PostHandler {
	return &PostHandler{createUsecase: usecase}
}

// CreatePostRequest は POST /posts の入力。
type CreatePostRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

// CreatePostResponse は作成結果を表す。
type CreatePostResponse struct {
	PostID string `json:"post_id"`
}

// CreatePost はリクエストを受けてユースケースを実行する。
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: messagePostInvalidRequest})
		return
	}
	if strings.TrimSpace(req.PostID) == "" || strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Message: messagePostInvalidRequest})
		return
	}

	out, err := h.createUsecase.Execute(c.Request.Context(), &postusecase.CreatePostInput{
		DarkPostID: req.PostID,
		Content:    req.Content,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, CreatePostResponse{PostID: out.DarkPostID})
}

func (h *PostHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, postusecase.ErrNilInput):
		c.JSON(http.StatusBadRequest, errorResponse{Message: messagePostInvalidRequest})
	case errors.Is(err, postdomain.ErrEmptyContent):
		c.JSON(http.StatusBadRequest, errorResponse{Message: messagePostInvalidRequest})
	case errors.Is(err, postusecase.ErrPostAlreadyExists),
		errors.Is(err, postusecase.ErrJobAlreadyScheduled):
		c.JSON(http.StatusConflict, errorResponse{Message: messagePostConflict})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Message: messageInternalError})
	}
}
