package batch_get_tweets

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets/mappers"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/list_items"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/timeline-service/internal/platform/get_timeline"
	"github.com/uala-challenge/timeline-service/kit"
)

type service struct {
	db   list_items.Service
	rd   get_timeline.Service
	log  log.Service
	conf Config
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		db:   d.DBRepository,
		rd:   d.RedisRepository,
		log:  d.Log,
		conf: d.Config,
	}
}

func (s service) Apply(ctx context.Context, user string) ([]kit.Tweet, error) {
	tweets := s.rd.Apply(ctx, fmt.Sprintf("timeline:%s", user))
	s.log.Info(ctx, fmt.Sprintf("tweets %s", tweets), nil)

	if len(tweets) == 0 {
		s.log.Info(ctx, "Timeline vacío, no se ejecuta BatchGetItem", nil)
		return make([]kit.Tweet, 0), nil
	}

	keys := GenerateTweetKeys(tweets)

	items, err := s.db.Apply(ctx, keys, s.conf.Table)
	if err != nil {
		return nil, err
	}

	var results []*kit.DynamoItem
	for _, item := range items {
		var dynamoItem kit.DynamoItem
		err := s.unmarshalDynamoItem(item, &dynamoItem)
		if err != nil {
			return nil, s.log.WrapError(err, "error al deserializar el item")
		}
		results = append(results, &dynamoItem)
	}

	return mappers.DynamoItemsToTweets(results), nil

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

func (s service) unmarshalDynamoItem(item map[string]types.AttributeValue, out interface{}) error {
	if item == nil {
		return fmt.Errorf("el item está vacío o es nil")
	}

	err := attributevalue.UnmarshalMap(item, out)
	if err != nil {
		return s.log.WrapError(err, "error al unmarshallar el item")
	}

	return nil
}
