package redis_timeline

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Apply(ctx context.Context, key string) []map[string]string
}

type Dependencies struct {
	Client *redis.Client
	Log    log.Service
}
