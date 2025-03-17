package mappers

import "github.com/uala-challenge/timeline-service/kit"

func DynamoItemToTweet(g *kit.DynamoItem) *kit.Tweet {
	return &kit.Tweet{
		UserID:  g.SK,
		TweetID: g.PK,
		Created: g.Created,
		Content: g.Content,
	}
}
