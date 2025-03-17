package update_timeline

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/timeline-service/kit"
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

func (s service) Accept(ctx context.Context, follower string, tweet kit.Tweet) error {
	_, err := s.client.HSet(ctx, tweet.TweetID, map[string]interface{}{
		"tweet_id":   tweet.TweetID,
		"user_id":    tweet.UserID,
		"created_at": tweet.Created,
	}).Result()

	if err != nil {
		return s.log.WrapError(err, "Error al guardar el item")
	}

	_, err = s.client.ZAdd(ctx, fmt.Sprintf("timeline:%s", follower), redis.Z{
		Score:  float64(tweet.Created),
		Member: tweet.TweetID,
	}).Result()

	if err != nil {
		return s.log.WrapError(err, "Error al guardar el item")
	}
	return nil
}
