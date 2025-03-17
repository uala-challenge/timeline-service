package patch_timeline

import (
	"net/http"

	"github.com/uala-challenge/timeline-service/internal/refresh_user_timeline"
)

type Service interface {
	Init(w http.ResponseWriter, r *http.Request)
}

type Dependencies struct {
	UseCaseRefresh refresh_user_timeline.Service
}
