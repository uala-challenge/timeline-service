package config

import (
	"github.com/uala-challenge/timeline-service/internal/platform/list_items"
)

type RepositoryConfig struct {
	Tweets list_items.Config `json:"tweets"`
}
