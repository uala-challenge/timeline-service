package patch_timeline

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cm "github.com/uala-challenge/timeline-service/internal/refresh_user_timeline/mock"
	"github.com/uala-challenge/timeline-service/kit"
)

func TestPatchTimeline_Success(t *testing.T) {
	mockUseCase := cm.NewService(t)

	mockUseCase.On("Accept", mock.Anything, "user-123", "user-456").Return(nil)

	service := NewService(Dependencies{
		UseCaseRefresh: mockUseCase,
	})

	requestBody := kit.Request{FollowerID: "user-456"}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PATCH", "/timeline/user-123", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("user_id", "user-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Timeline updated successfully", response["message"])

	mockUseCase.AssertExpectations(t)
}

func TestPatchTimeline_Error(t *testing.T) {
	mockUseCase := cm.NewService(t)

	mockUseCase.On("Accept", mock.Anything, "user-123", "user-456").Return(errors.New("database error"))

	service := NewService(Dependencies{
		UseCaseRefresh: mockUseCase,
	})

	requestBody := kit.Request{FollowerID: "user-456"}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PATCH", "/timeline/user-123", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("user_id", "user-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error tu update timeline")

	mockUseCase.AssertExpectations(t)
}
