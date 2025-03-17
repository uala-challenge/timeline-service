package update_timeline

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/timeline-service/kit"
)

type Service interface {
	Accept(ctx context.Context, user string, tweet kit.Tweet) error
}

type Dependencies struct {
	Client *redis.Client
	Log    log.Service
}
