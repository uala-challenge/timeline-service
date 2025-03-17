package config

import (
	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"
	"github.com/uala-challenge/timeline-service/internal/refresh_user_timeline"
)

type UsesCasesConfig struct {
	Tweets  batch_get_tweets.Config      `json:"tweets"`
	Refresh refresh_user_timeline.Config `json:"refresh"`
}
