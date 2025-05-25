package app

import (
	"context"
	"gophernet/pkg/cache"
	"gophernet/pkg/repo"
)

type IGopherApp interface {
	GetGopher() string
	StartScheduler(ctx context.Context)
	StopScheduler()
}

type GopherApp struct {
	repo      repo.IBurrowRepository
	scheduler *Scheduler
}

func NewGopherApp(repo repo.IBurrowRepository, stats cache.IBurrowStats) *GopherApp {
	return &GopherApp{
		repo:      repo,
		scheduler: NewScheduler(repo, stats),
	}
}

func (g *GopherApp) GetGopher() string {
	return "gopher"
}

func (g *GopherApp) StartScheduler(ctx context.Context) {
	g.scheduler.Start(ctx)
}

func (g *GopherApp) StopScheduler() {
	g.scheduler.Stop()
}
