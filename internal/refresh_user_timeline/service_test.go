package refresh_user_timeline

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	query_mock "github.com/uala-challenge/simple-toolkit/pkg/platform/db/query/mock"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
	update_mock "github.com/uala-challenge/timeline-service/internal/platform/update_timeline/mock"
)

func TestRefreshUserTimeline_Success(t *testing.T) {
	mockDB := query_mock.NewService(t)
	mockRedis := update_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	userID := "user:123"
	followerID := "user:456"
	expectedItems := []map[string]types.AttributeValue{
		{
			"PK":        &types.AttributeValueMemberS{Value: "tweet:1"},
			"SK":        &types.AttributeValueMemberS{Value: "user:123"},
			"createdAt": &types.AttributeValueMemberN{Value: "1710087745"},
		},
	}

	mockDB.On("Apply", mock.Anything, mock.Anything).Return(expectedItems, nil)
	mockRedis.On("Accept", mock.Anything, followerID, mock.Anything).Return(nil)

	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
	})

	err := service.Accept(context.TODO(), userID, followerID)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestRefreshUserTimeline_NoTweets(t *testing.T) {
	mockDB := query_mock.NewService(t)
	mockRedis := update_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	userID := "user:123"
	followerID := "user:456"
	mockDB.On("Apply", mock.Anything, mock.Anything).Return([]map[string]types.AttributeValue{}, nil)

	mockLog.On("Debug", mock.Anything, mock.Anything, mock.Anything).Return()

	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
	})

	err := service.Accept(context.TODO(), userID, followerID)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestRefreshUserTimeline_DynamoError(t *testing.T) {
	mockDB := query_mock.NewService(t)
	mockRedis := update_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	userID := "user:123"
	followerID := "user:456"
	expectedErr := fmt.Errorf("error en DynamoDB")
	mockDB.On("Apply", mock.Anything, mock.Anything).Return(nil, expectedErr)

	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
	})

	err := service.Accept(context.TODO(), userID, followerID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

func TestRefreshUserTimeline_UnmarshalError(t *testing.T) {
	mockDB := query_mock.NewService(t)
	mockRedis := update_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	userID := "user:123"
	followerID := "user:456"
	invalidItems := []map[string]types.AttributeValue{
		{
			"PK": &types.AttributeValueMemberB{Value: []byte("invalidData")}, // Tipo incorrecto
		},
	}

	mockDB.On("Apply", mock.Anything, mock.Anything).Return(invalidItems, nil)

	mockLog.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
	})

	err := service.Accept(context.TODO(), userID, followerID)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestRefreshUserTimeline_RedisError(t *testing.T) {
	mockDB := query_mock.NewService(t)
	mockRedis := update_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	userID := "user:123"
	followerID := "user:456"
	expectedItems := []map[string]types.AttributeValue{
		{
			"PK":        &types.AttributeValueMemberS{Value: "tweet:1"},
			"SK":        &types.AttributeValueMemberS{Value: "user:123"},
			"createdAt": &types.AttributeValueMemberN{Value: "1710087745"},
		},
	}

	expectedErr := fmt.Errorf("error en Redis")
	mockDB.On("Apply", mock.Anything, mock.Anything).Return(expectedItems, nil)
	mockRedis.On("Accept", mock.Anything, followerID, mock.Anything).Return(expectedErr)
	mockLog.On("Error", mock.Anything, expectedErr, mock.Anything, mock.Anything).Return()

	service := NewService(Dependencies{
		DBRepository:    mockDB,
		RedisRepository: mockRedis,
		Log:             mockLog,
	})

	err := service.Accept(context.TODO(), userID, followerID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "falló la actualización de 1 tweets en Redis")
	mockDB.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}
