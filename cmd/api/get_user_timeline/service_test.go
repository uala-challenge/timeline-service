package get_user_timeline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cm "github.com/uala-challenge/timeline-service/internal/batch_get_tweets/mock"
	"github.com/uala-challenge/timeline-service/kit"
)

func TestGetTimeline_Success(t *testing.T) {
	mockUseCase := cm.NewService(t)

	mockTweets := []kit.Tweet{
		{TweetID: "tweet-1", UserID: "user-123", Content: "Hello World!"},
		{TweetID: "tweet-2", UserID: "user-123", Content: "Another tweet"},
	}

	mockUseCase.On("Apply", mock.Anything, "user-123").Return(mockTweets, nil)

	service := NewService(Dependencies{
		UseCaseRetrieveTweet: mockUseCase,
	})

	req := httptest.NewRequest("GET", "/timeline/user-123", nil)
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("user_id", "user-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response []kit.Tweet
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockTweets, response)

	mockUseCase.AssertExpectations(t)
}

func TestGetTimeline_NotFound(t *testing.T) {
	mockUseCase := cm.NewService(t)

	mockUseCase.On("Apply", mock.Anything, "user-123").Return([]kit.Tweet{}, nil)

	service := NewService(Dependencies{
		UseCaseRetrieveTweet: mockUseCase,
	})

	req := httptest.NewRequest("GET", "/timeline/user-123", nil)
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("user_id", "user-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "tweet not found")

	mockUseCase.AssertExpectations(t)
}

func TestGetTimeline_Error(t *testing.T) {
	mockUseCase := cm.NewService(t)

	mockUseCase.On("Apply", mock.Anything, "user-123").Return(nil, errors.New("database error"))

	service := NewService(Dependencies{
		UseCaseRetrieveTweet: mockUseCase,
	})

	req := httptest.NewRequest("GET", "/timeline/user-123", nil)
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("user_id", "user-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error to get tweet")

	mockUseCase.AssertExpectations(t)
}
