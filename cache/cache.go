package cache

import (
	"log/slog"
	"time"

	"github.com/deparr/api/model"
)

type githubCache struct {
	pinned []model.Repository
	recent []model.Repository
	updated time.Time
}

type cache struct {
	github githubCache
}

var globalCache cache
var cacheTickers struct{
	github *time.Ticker
}
var toggleChans struct {
	github chan bool
}

func (c *cache) setGithub(which string, repos []model.Repository, updated time.Time) {
	if repos == nil {
		return
	}

	if which == "pinned" {
		c.github.pinned = repos
	} else {
		c.github.recent = repos
	}

	slog.Info("updated cache(gh)", "which", which)

	c.github.updated = updated
}

func InitCache() {
	go initRepoCache()
}

func GetGithubPinned() []model.Repository {
	if globalCache.github.pinned == nil {
		slog.Info("cold cache(gh.pinned)")
		updateRepoCache(time.Now())
	}

	return globalCache.github.pinned
}

func GetGithubRecent() []model.Repository {
	if globalCache.github.recent == nil {
		slog.Info("cold cache(gh.recent)")
		updateRepoCache(time.Now())
	}

	return globalCache.github.recent
}
