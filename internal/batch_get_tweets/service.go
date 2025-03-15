package batch_get_tweets

import (
	"context"
	"fmt"
	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets/mappers"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/timeline-service/internal/platform/list_items"
	"github.com/uala-challenge/timeline-service/internal/platform/redis_timeline"
	"github.com/uala-challenge/timeline-service/kit"
)

type service struct {
	db  list_items.Service
	rd  redis_timeline.Service
	log log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		db:  d.DBRepository,
		rd:  d.RedisRepository,
		log: d.Log,
	}
}

func (s service) Apply(ctx context.Context, user string) ([]kit.Tweet, error) {
	tweets := s.rd.Apply(ctx, fmt.Sprintf("timeline:%s", user))
	s.log.Info(ctx, fmt.Sprintf("tweets %s", tweets), nil)

	keys := GenerateTweetKeys(tweets)

	items, err := s.db.Apply(ctx, keys)
	if err != nil {
		return nil, err
	}

	return mappers.DynamoItemsToTweets(items), nil

}

func GenerateTweetKeys(tweets []map[string]string) []map[string]types.AttributeValue {
	var keys []map[string]types.AttributeValue

	for _, tweet := range tweets {
		key := map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: tweet["tweet_id"]},
			"SK": &types.AttributeValueMemberS{Value: tweet["user_id"]},
		}
		keys = append(keys, key)
	}

	return keys
}
