package main

import (
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/list_items"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_engine"
	"github.com/uala-challenge/timeline-service/cmd/api/get_timeline"
	"github.com/uala-challenge/timeline-service/internal/batch_get_tweets"
	"github.com/uala-challenge/timeline-service/internal/platform/redis_timeline"
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
	a.engine.repositories.GetTimeline = redis_timeline.NewService(redis_timeline.Dependencies{
		Client: a.engine.simplify.RedisClient,
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
	return a
}

func (a AppBuilder) InitHandlers() app_builder.Builder {
	a.engine.handlers.GetTimeline = get_timeline.
		NewService(get_timeline.Dependencies{
			UseCaseRetrieveTweet: a.engine.useCases.BatchGetTweets})
	return a
}

func (a AppBuilder) InitRoutes() app_builder.Builder {
	a.engine.simplify.App.Router.Get("/timeline/{user_id}", a.engine.handlers.GetTimeline.Init)
	return a
}

func (a AppBuilder) Build() app_builder.App {
	return a.engine
}
