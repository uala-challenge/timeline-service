package refresh_user_timeline

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/query"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/timeline-service/internal/platform/update_timeline"
	"github.com/uala-challenge/timeline-service/internal/refresh_user_timeline/mappers"
	"github.com/uala-challenge/timeline-service/kit"
)

type service struct {
	db   query.Service
	rd   update_timeline.Service
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

func (s service) Accept(ctx context.Context, user, follower string) error {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("UalaChallenge"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: user},
		},
	}

	values, err := s.db.Apply(ctx, input)
	if err != nil {
		return err
	}

	if len(values) == 0 {
		s.log.Debug(ctx, map[string]interface{}{"tweets_count": len(values),
			"mensaje": "No se encontraron tweets para actualizar en el timeline"})
		return nil
	}

	var failedTweets []kit.DynamoItem

	for _, item := range values {
		var dynamoItem kit.DynamoItem
		if err := attributevalue.UnmarshalMap(item, &dynamoItem); err != nil {
			s.log.Error(ctx, err, "Error al deserializar el item", map[string]interface{}{
				"error": err.Error(),
				"item":  item,
			})
			continue
		}

		if err := s.rd.Accept(ctx, follower, *mappers.DynamoItemToTweet(&dynamoItem)); err != nil {
			s.log.Error(ctx, err, "Error al actualizar timeline en Redis", map[string]interface{}{
				"error": err.Error(),
				"user":  user,
				"tweet": dynamoItem,
			})
			failedTweets = append(failedTweets, dynamoItem)
		}
	}

	if len(failedTweets) > 0 {
		return fmt.Errorf("falló la actualización de %d tweets en Redis", len(failedTweets))
	}
	return nil
}
