package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	postdomain "backend/internal/domain/post"
	postusecase "backend/internal/usecase/post"

	"github.com/gin-gonic/gin"
)

func TestPostHandler_CreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		stub := &stubCreatePostUsecase{
			output: &postusecase.CreatePostOutput{DarkPostID: "dark-1"},
		}
		handler := NewPostHandler(stub)
		router := gin.New()
		router.POST("/posts", handler.CreatePost)

		rec := httptest.NewRecorder()
		reqBody := bytes.NewBufferString(`{"post_id":"dark-1","content":"hello"}`)
		req := httptest.NewRequest(http.MethodPost, "/posts", reqBody)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		if stub.received.DarkPostID != "dark-1" || stub.received.Content != "hello" {
			t.Fatalf("unexpected input passed to usecase: %+v", stub.received)
		}

		var resp CreatePostResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.PostID != "dark-1" {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		handler := NewPostHandler(&stubCreatePostUsecase{})
		rec, resp := performPostRequest(handler, `{"post_id":`)
		expectStatusAndMessage(t, rec, resp, http.StatusBadRequest, messagePostInvalidRequest)
	})

	t.Run("empty fields", func(t *testing.T) {
		handler := NewPostHandler(&stubCreatePostUsecase{})
		rec, resp := performPostRequest(handler, `{"post_id":"","content":""}`)
		expectStatusAndMessage(t, rec, resp, http.StatusBadRequest, messagePostInvalidRequest)
	})

	t.Run("post already exists", func(t *testing.T) {
		handler := NewPostHandler(&stubCreatePostUsecase{
			err: postusecase.ErrPostAlreadyExists,
		})
		rec, resp := performPostRequest(handler, `{"post_id":"dup","content":"hello"}`)
		expectStatusAndMessage(t, rec, resp, http.StatusConflict, messagePostConflict)
	})

	t.Run("domain validation error", func(t *testing.T) {
		handler := NewPostHandler(&stubCreatePostUsecase{
			err: postdomain.ErrEmptyContent,
		})
		rec, resp := performPostRequest(handler, `{"post_id":"dark","content":"hello"}`)
		expectStatusAndMessage(t, rec, resp, http.StatusBadRequest, messagePostInvalidRequest)
	})

	t.Run("nil input error", func(t *testing.T) {
		handler := NewPostHandler(&stubCreatePostUsecase{
			err: postusecase.ErrNilInput,
		})
		rec, resp := performPostRequest(handler, `{"post_id":"dark","content":"hello"}`)
		expectStatusAndMessage(t, rec, resp, http.StatusBadRequest, messagePostInvalidRequest)
	})

	t.Run("job already scheduled", func(t *testing.T) {
		handler := NewPostHandler(&stubCreatePostUsecase{
			err: postusecase.ErrJobAlreadyScheduled,
		})
		rec, resp := performPostRequest(handler, `{"post_id":"dark","content":"hello"}`)
		expectStatusAndMessage(t, rec, resp, http.StatusConflict, messagePostConflict)
	})

	t.Run("internal error", func(t *testing.T) {
		handler := NewPostHandler(&stubCreatePostUsecase{
			err: errors.New("boom"),
		})
		rec, resp := performPostRequest(handler, `{"post_id":"dark","content":"hello"}`)
		expectStatusAndMessage(t, rec, resp, http.StatusInternalServerError, messageInternalError)
	})
}

func performPostRequest(handler *PostHandler, body string) (*httptest.ResponseRecorder, []byte) {
	router := gin.New()
	router.POST("/posts", handler.CreatePost)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	return rec, rec.Body.Bytes()
}

func expectStatusAndMessage(t *testing.T, rec *httptest.ResponseRecorder, body []byte, status int, message string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("expected status %d, got %d", status, rec.Code)
	}
	var resp errorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != message {
		t.Fatalf("expected message %q, got %q", message, resp.Message)
	}
}

type stubCreatePostUsecase struct {
	output   *postusecase.CreatePostOutput
	err      error
	received *postusecase.CreatePostInput
}

func (s *stubCreatePostUsecase) Execute(ctx context.Context, in *postusecase.CreatePostInput) (*postusecase.CreatePostOutput, error) {
	s.received = in
	if s.err != nil {
		return nil, s.err
	}
	if s.output != nil {
		return s.output, nil
	}
	return &postusecase.CreatePostOutput{DarkPostID: in.DarkPostID}, nil
}
