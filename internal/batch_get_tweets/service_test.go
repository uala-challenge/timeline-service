package batch_get_tweets

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	db_mock "github.com/uala-challenge/simple-toolkit/pkg/platform/db/list_items/mock"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
	redis_mock "github.com/uala-challenge/timeline-service/internal/platform/redis_timeline/mock"
)

func TestBatchGetTweets_EmptyTimeline(t *testing.T) {
	mockRedis := redis_mock.NewService(t)
	mockDB := db_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	mockRedis.On("Apply", mock.Anything, "timeline:user-123").Return([]map[string]string{})

	mockLog.On("Info", mock.Anything, "tweets []", mock.Anything).Return()
	mockLog.On("Info", mock.Anything, "Timeline vacío, no se ejecuta BatchGetItem", mock.Anything).Return()

	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
		Config:          Config{Table: "tweets_table"},
	})

	tweets, err := service.Apply(context.TODO(), "user-123")

	assert.NoError(t, err)
	assert.Empty(t, tweets)

	mockRedis.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestBatchGetTweets_Success(t *testing.T) {
	mockRedis := redis_mock.NewService(t)
	mockDB := db_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	mockRedis.On("Apply", mock.Anything, "timeline:user-123").Return([]map[string]string{
		{"tweet_id": "tweet-1", "user_id": "user-123"},
		{"tweet_id": "tweet-2", "user_id": "user-123"},
	})

	mockKeys := []map[string]types.AttributeValue{
		{"PK": &types.AttributeValueMemberS{Value: "tweet-1"}, "SK": &types.AttributeValueMemberS{Value: "user-123"}},
		{"PK": &types.AttributeValueMemberS{Value: "tweet-2"}, "SK": &types.AttributeValueMemberS{Value: "user-123"}},
	}

	mockDB.On("Apply", mock.Anything, mockKeys, "tweets_table").Return([]map[string]types.AttributeValue{
		{"PK": &types.AttributeValueMemberS{Value: "tweet-1"}, "SK": &types.AttributeValueMemberS{Value: "user-123"}, "Content": &types.AttributeValueMemberS{Value: "Hello World"}},
		{"PK": &types.AttributeValueMemberS{Value: "tweet-2"}, "SK": &types.AttributeValueMemberS{Value: "user-123"}, "Content": &types.AttributeValueMemberS{Value: "Go is awesome"}},
	}, nil)

	mockLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
		Config:          Config{Table: "tweets_table"},
	})

	tweets, err := service.Apply(context.TODO(), "user-123")

	assert.NoError(t, err)
	assert.Len(t, tweets, 2)
	assert.Equal(t, "Hello World", tweets[0].Content)
	assert.Equal(t, "Go is awesome", tweets[1].Content)

	mockRedis.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestBatchGetTweets_DBError(t *testing.T) {
	mockRedis := redis_mock.NewService(t)
	mockDB := db_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	mockRedis.On("Apply", mock.Anything, "timeline:user-123").Return([]map[string]string{
		{"tweet_id": "tweet-1", "user_id": "user-123"},
		{"tweet_id": "tweet-2", "user_id": "user-123"},
	})

	mockKeys := []map[string]types.AttributeValue{
		{"PK": &types.AttributeValueMemberS{Value: "tweet-1"}, "SK": &types.AttributeValueMemberS{Value: "user-123"}},
		{"PK": &types.AttributeValueMemberS{Value: "tweet-2"}, "SK": &types.AttributeValueMemberS{Value: "user-123"}},
	}

	mockDB.On("Apply", mock.Anything, mockKeys, "tweets_table").Return(nil, assert.AnError)
	mockLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
		Config:          Config{Table: "tweets_table"},
	})

	tweets, err := service.Apply(context.TODO(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, tweets)
}
