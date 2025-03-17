package main

import (
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/list_items"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/query"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_engine"
	"github.com/uala-challenge/timeline-service/cmd/api/get_user_timeline"
	"github.com/uala-challenge/timeline-service/cmd/api/patch_timeline"
	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"
	"github.com/uala-challenge/timeline-service/internal/platform/get_timeline"
	"github.com/uala-challenge/timeline-service/internal/platform/update_timeline"
	"github.com/uala-challenge/timeline-service/internal/refresh_user_timeline"

	"github.com/uala-challenge/timeline-service/kit/config"
)

type engine struct {
	simplify       app_engine.Engine
	repositories   repositories
	useCases       useCases
	handlers       handlers
	useCasesConfig config.UsesCasesConfig
}

type AppBuilder struct {
	engine *engine
}

var _ app_builder.Builder = (*AppBuilder)(nil)

func NewAppBuilder() *AppBuilder {
	a := *app_engine.NewApp()
	return &AppBuilder{
		engine: &engine{
			simplify: a,
		},
	}
}

func (a engine) Run() error {
	return a.simplify.App.Run()
}

func (a AppBuilder) LoadConfig() app_builder.Builder {
	a.engine.useCasesConfig = app_engine.GetConfig[config.UsesCasesConfig](a.engine.simplify.UsesCasesConfig)
	return a
}

func (a AppBuilder) InitRepositories() app_builder.Builder {
	a.engine.repositories.GetTweets = list_items.NewService(list_items.Dependencies{
		Client: a.engine.simplify.DynamoDBClient,
		Log:    a.engine.simplify.Log,
	})
	a.engine.repositories.GetTimeline = get_timeline.NewService(get_timeline.Dependencies{
		Client: a.engine.simplify.RedisClient,
		Log:    a.engine.simplify.Log,
	})
	a.engine.repositories.UpdateTimeline = update_timeline.NewService(update_timeline.Dependencies{
		Client: a.engine.simplify.RedisClient,
		Log:    a.engine.simplify.Log,
	})
	a.engine.repositories.GetUserTweets = query.NewService(query.Dependencies{
		Client: a.engine.simplify.DynamoDBClient,
		Log:    a.engine.simplify.Log,
	})
	return a
}

func (a AppBuilder) InitUseCases() app_builder.Builder {
	a.engine.useCases.BatchGetTweets = batch_get_tweets.NewService(batch_get_tweets.Dependencies{
		DBRepository:    a.engine.repositories.GetTweets,
		RedisRepository: a.engine.repositories.GetTimeline,
		Log:             a.engine.simplify.Log,
		Config:          a.engine.useCasesConfig.Tweets,
	})
	a.engine.useCases.RefreshTimeLine = refresh_user_timeline.NewService(refresh_user_timeline.Dependencies{
		DBRepository:    a.engine.repositories.GetUserTweets,
		RedisRepository: a.engine.repositories.UpdateTimeline,
		Log:             a.engine.simplify.Log,
		Config:          a.engine.useCasesConfig.Refresh,
	})
	return a
}

func (a AppBuilder) InitHandlers() app_builder.Builder {
	a.engine.handlers.GetUserTimeline = get_user_timeline.
		NewService(get_user_timeline.Dependencies{
			UseCaseRetrieveTweet: a.engine.useCases.BatchGetTweets})
	a.engine.handlers.PatchTimeline = patch_timeline.NewService(patch_timeline.Dependencies{
		UseCaseRefresh: a.engine.useCases.RefreshTimeLine})
	return a
}

func (a AppBuilder) InitRoutes() app_builder.Builder {
	a.engine.simplify.App.Router.Get("/timeline/{user_id}", a.engine.handlers.GetUserTimeline.Init)
	a.engine.simplify.App.Router.Patch("/timeline/{user_id}", a.engine.handlers.PatchTimeline.Init)
	return a
}

func (a AppBuilder) Build() app_builder.App {
	return a.engine
}
