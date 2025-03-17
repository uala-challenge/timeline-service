package main

import (
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/list_items"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/query"
	"github.com/uala-challenge/timeline-service/cmd/api/get_user_timeline"
	"github.com/uala-challenge/timeline-service/cmd/api/patch_timeline"
	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"
	"github.com/uala-challenge/timeline-service/internal/platform/get_timeline"
	"github.com/uala-challenge/timeline-service/internal/platform/update_timeline"
	"github.com/uala-challenge/timeline-service/internal/refresh_user_timeline"
)

type repositories struct {
	GetTweets      list_items.Service
	GetTimeline    get_timeline.Service
	UpdateTimeline update_timeline.Service
	GetUserTweets  query.Service
}

type useCases struct {
	BatchGetTweets  batch_get_tweets.Service
	RefreshTimeLine refresh_user_timeline.Service
}

type handlers struct {
	GetUserTimeline get_user_timeline.Service
	PatchTimeline   patch_timeline.Service
}
