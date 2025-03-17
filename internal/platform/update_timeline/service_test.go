package update_timeline

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
	"github.com/uala-challenge/timeline-service/kit"
)

func TestUpdateTimeline_Success(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet:1",
		UserID:  "user:123",
		Created: 1710087745,
	}

	mk.ExpectHSet("tweet:1", map[string]interface{}{
		"tweet_id":   tweet.TweetID,
		"user_id":    tweet.UserID,
		"created_at": tweet.Created,
	}).SetVal(1)

	mk.ExpectZAdd("timeline:user:123", redis.Z{
		Score:  float64(tweet.Created),
		Member: tweet.TweetID,
	}).SetVal(1)

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	err := service.Accept(context.TODO(), "user:123", tweet)

	assert.NoError(t, err)

	err = mk.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateTimeline_ErrorHSet(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet:1",
		UserID:  "user:123",
		Created: 1710087745,
	}

	expectedErr := fmt.Errorf("error en HSet")
	mk.ExpectHSet("tweet:1", map[string]interface{}{
		"tweet_id":   tweet.TweetID,
		"user_id":    tweet.UserID,
		"created_at": tweet.Created,
	}).SetErr(expectedErr)

	mockLog.On("WrapError", expectedErr, "Error al guardar el item").Return(expectedErr)

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	err := service.Accept(context.TODO(), "user:123", tweet)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	err = mk.ExpectationsWereMet()
	assert.NoError(t, err)
	mockLog.AssertExpectations(t)
}

func TestUpdateTimeline_ErrorZAdd(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweet := kit.Tweet{
		TweetID: "tweet:1",
		UserID:  "user:123",
		Created: 1710087745,
	}

	mk.ExpectHSet("tweet:1", map[string]interface{}{
		"tweet_id":   tweet.TweetID,
		"user_id":    tweet.UserID,
		"created_at": tweet.Created,
	}).SetVal(1)

	expectedErr := fmt.Errorf("error en ZAdd")
	mk.ExpectZAdd("timeline:user:123", redis.Z{
		Score:  float64(tweet.Created),
		Member: tweet.TweetID,
	}).SetErr(expectedErr)

	mockLog.On("WrapError", expectedErr, "Error al guardar el item").Return(expectedErr)

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	err := service.Accept(context.TODO(), "user:123", tweet)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	err = mk.ExpectationsWereMet()
	assert.NoError(t, err)
	mockLog.AssertExpectations(t)
}
