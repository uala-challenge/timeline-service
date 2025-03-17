package config

import (
	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"
)

type UsesCasesConfig struct {
	Tweets batch_get_tweets.Config `json:"tweets"`
}
