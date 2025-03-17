package main

import (
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/list_items"
	"github.com/uala-challenge/timeline-service/cmd/api/get_timeline"
	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"
	"github.com/uala-challenge/timeline-service/internal/platform/redis_timeline"
)

type repositories struct {
	GetTweets   list_items.Service
	GetTimeline redis_timeline.Service
}

type useCases struct {
	BatchGetTweets batch_get_tweets.Service
}

type handlers struct {
	GetTimeline get_timeline.Service
}
