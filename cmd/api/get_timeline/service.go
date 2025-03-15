package get_timeline

import (
	"encoding/json"
	"net/http"

	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"

	"github.com/go-chi/chi/v5"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/error_handler"
)

type service struct {
	useCase batch_get_tweets.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		useCase: d.UseCaseRetrieveTweet,
	}
}

func (s service) Init(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	tweet, err := s.useCase.Apply(r.Context(), userID)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("error to get tweet", err.Error(), err, http.StatusInternalServerError), w)
		return
	}

	rsp, err := json.Marshal(tweet)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("error to marshal tweet", err.Error(), err, http.StatusInternalServerError), w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(rsp)
}
