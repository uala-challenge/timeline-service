package batch_get_tweets

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/list_items"
	"github.com/uala-challenge/timeline-service/internal/platform/get_timeline"
	"github.com/uala-challenge/timeline-service/kit"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Apply(ctx context.Context, user string) ([]kit.Tweet, error)
}

type Dependencies struct {
	DBRepository    list_items.Service
	RedisRepository get_timeline.Service
	Log             log.Service
	Config          Config
}

type Config struct {
	Table string `json:"table"`
}
