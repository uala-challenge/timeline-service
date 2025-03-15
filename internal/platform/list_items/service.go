package list_items

import (
	"context"
	"fmt"

	"github.com/uala-challenge/timeline-service/kit"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *dynamodb.Client
	log    log.Service
	config Config
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	return &service{
		client: d.Client,
		log:    d.Log,
		config: d.Config,
	}
}

func (s service) Apply(ctx context.Context, items []map[string]types.AttributeValue) ([]*kit.DynamoItem, error) {
	batchInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			s.config.Table: {Keys: items},
		},
	}

	output, err := s.client.BatchGetItem(ctx, batchInput)
	if err != nil {
		return nil, s.log.WrapError(err, "error al ejecutar BatchGetItem")
	}

	var results []*kit.DynamoItem
	for _, item := range output.Responses[s.config.Table] {
		var dynamoItem kit.DynamoItem
		err := s.unmarshalDynamoItem(item, &dynamoItem)
		if err != nil {
			return nil, s.log.WrapError(err, "error al deserializar el item")
		}
		results = append(results, &dynamoItem)
	}

	return results, nil
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
