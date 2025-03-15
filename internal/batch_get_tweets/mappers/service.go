package mappers

import "github.com/uala-challenge/timeline-service/kit"

func DynamoItemsToTweets(g []*kit.DynamoItem) []kit.Tweet {
	var tweets []kit.Tweet

	for _, item := range g {
		tweet := kit.Tweet{
			UserID:  item.SK,
			TweetID: item.PK,
			Created: item.Created,
			Content: item.Content,
		}
		tweets = append(tweets, tweet)
	}

	return tweets
}
