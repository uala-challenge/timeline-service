package get_timeline

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestRedisTimeline_Success(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweetIDs := []string{"tweet:1", "tweet:2"}
	tweetData1 := map[string]string{"id": "tweet:1", "content": "Primer tweet"}
	tweetData2 := map[string]string{"id": "tweet:2", "content": "Segundo tweet"}

	// Simular que Redis devuelve tweets en la línea de tiempo
	mk.ExpectZRevRange("user:123", 0, -1).SetVal(tweetIDs)
	mk.ExpectHGetAll("tweet:1").SetVal(tweetData1)
	mk.ExpectHGetAll("tweet:2").SetVal(tweetData2)

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	tweets := service.Apply(context.TODO(), "user:123")

	assert.Len(t, tweets, 2)
	assert.Equal(t, "tweet:1", tweets[0]["id"])
	assert.Equal(t, "Primer tweet", tweets[0]["content"])
	assert.Equal(t, "tweet:2", tweets[1]["id"])
	assert.Equal(t, "Segundo tweet", tweets[1]["content"])

	err := mk.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisTimeline_ErrorZRevRange(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	expectedErr := fmt.Errorf("error en Redis")
	mk.ExpectZRevRange("user:123", 0, -1).SetErr(expectedErr)
	mockLog.On("Error", mock.Anything, expectedErr, "Error obteniendo tweets del timeline", mock.Anything).Return()
	mockLog.On("Info", mock.Anything, "No tweets del timeline", mock.Anything).Return()

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	tweets := service.Apply(context.TODO(), "user:123")

	assert.Nil(t, tweets)

	err := mk.ExpectationsWereMet()
	assert.NoError(t, err)
	mockLog.AssertExpectations(t)
}

func TestRedisTimeline_ErrorHGetAll(t *testing.T) {
	mockRedis, mk := redismock.NewClientMock()
	mockLog := log_mock.NewService(t)

	tweetIDs := []string{"tweet:1", "tweet:2"}
	tweetData1 := map[string]string{"id": "tweet:1", "content": "Primer tweet"}
	expectedErr := fmt.Errorf("error en Redis")

	mk.ExpectZRevRange("user:123", 0, -1).SetVal(tweetIDs)
	mk.ExpectHGetAll("tweet:1").SetVal(tweetData1)
	mk.ExpectHGetAll("tweet:2").SetErr(expectedErr)

	mockLog.On("Error", mock.Anything, expectedErr, mock.Anything, mock.Anything).Return()

	service := NewService(Dependencies{
		Client: mockRedis,
		Log:    mockLog,
	})

	tweets := service.Apply(context.TODO(), "user:123")

	assert.Len(t, tweets, 1)
	assert.Equal(t, "tweet:1", tweets[0]["id"])
	assert.Equal(t, "Primer tweet", tweets[0]["content"])

	err := mk.ExpectationsWereMet()
	assert.NoError(t, err)
	mockLog.AssertExpectations(t)
}
