package get_user_timeline

import (
	"net/http"

	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"
)

type Service interface {
	Init(w http.ResponseWriter, r *http.Request)
}

type Dependencies struct {
	UseCaseRetrieveTweet batch_get_tweets.Service
}
