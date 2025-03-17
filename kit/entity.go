package kit

type DynamoItem struct {
	PK      string `dynamodbav:"PK"`
	SK      string `dynamodbav:"SK"`
	GSI1PK  string `dynamodbav:"GSI1PK"`
	GSI1SK  string `dynamodbav:"GSI1SK"`
	Content string `dynamodbav:"content"`
	Created int64  `dynamodbav:"created"`
}

type Request struct {
	FollowerID string `json:"follower_id" validate:"required"`
}

type Tweet struct {
	UserID  string `json:"user_id"`
	TweetID string `json:"tweet_id"`
	Created int64  `json:"created"`
	Content string `json:"content"`
}
