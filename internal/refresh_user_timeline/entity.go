package refresh_user_timeline

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/query"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/timeline-service/internal/platform/update_timeline"
)

type Service interface {
	Accept(ctx context.Context, user, follower string) error
}

type Dependencies struct {
	DBRepository    query.Service
	RedisRepository update_timeline.Service
	Log             log.Service
	Config          Config
}

type Config struct {
	Table string `json:"table"`
}
