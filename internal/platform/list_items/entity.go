package list_items

import (
	"context"

	"github.com/uala-challenge/timeline-service/kit"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Apply(ctx context.Context, items []map[string]types.AttributeValue) ([]*kit.DynamoItem, error)
}

type Dependencies struct {
	Client *dynamodb.Client
	Config Config
	Log    log.Service
}

type Config struct {
	Table string `json:"table"`
}
