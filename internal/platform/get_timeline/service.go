package get_timeline

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *redis.Client
	log    log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	return &service{
		client: d.Client,
		log:    d.Log,
	}
}

func (s service) Apply(ctx context.Context, key string) []map[string]string {
	tweetIDs, err := s.client.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		s.log.Error(ctx, err, "Error obteniendo tweets del timeline", nil)
	}

	if len(tweetIDs) == 0 {
		s.log.Info(ctx, "No tweets del timeline", nil)
		return nil
	}

	pipe := s.client.Pipeline()
	cs := make([]*redis.MapStringStringCmd, len(tweetIDs))

	for i, tweetID := range tweetIDs {
		cs[i] = pipe.HGetAll(ctx, tweetID)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		s.log.Error(ctx, err, fmt.Sprintf("tweets %s", tweetIDs), nil)
	}

	tweets := []map[string]string{}

	for _, cmd := range cs {
		tweetData, err := cmd.Result()
		if err == nil && len(tweetData) > 0 {
			tweets = append(tweets, tweetData)
		}
	}

	return tweets

}
