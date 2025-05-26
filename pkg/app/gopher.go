package app

import (
	"context"
	"fmt"
	"gophernet/pkg/cache"
	"gophernet/pkg/db/ent"
	"gophernet/pkg/repo"
)

type IGopherApp interface {
	GetGopher() string
	StartScheduler(ctx context.Context)
	StopScheduler()
	RentBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error)
	ReleaseBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error)
	GetBurrowStatus(ctx context.Context) ([]*ent.Burrow, error)
}

type GopherApp struct {
	repo      repo.IBurrowRepository
	stats     cache.IBurrowStats
	scheduler *Scheduler
}

func NewGopherApp(repo repo.IBurrowRepository, stats cache.IBurrowStats) *GopherApp {
	return &GopherApp{
		repo:      repo,
		stats:     stats,
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

func (g *GopherApp) RentBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error) {
	burrow, err := g.repo.GetBurrowByID(ctx, burrowID)
	if err != nil {
		return nil, fmt.Errorf("burrow with ID %d not found: %w", burrowID, err)
	}

	if burrow.IsOccupied {
		return nil, fmt.Errorf("burrow %d is already occupied", burrowID)
	}

	if err := g.repo.UpdateBurrowOccupancy(ctx, burrowID, true); err != nil {
		return nil, fmt.Errorf("failed to update burrow %d occupancy: %w", burrowID, err)
	}
	g.stats.SubtractAvailableBurrow()
	burrow.IsOccupied = true
	return burrow, nil
}

func (g *GopherApp) ReleaseBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error) {
	burrow, err := g.repo.GetBurrowByID(ctx, burrowID)
	if err != nil {
		return nil, fmt.Errorf("burrow with ID %d not found: %w", burrowID, err)
	}

	if !burrow.IsOccupied {
		return nil, fmt.Errorf("burrow %d is not occupied", burrowID)
	}

	if err := g.repo.UpdateBurrowOccupancy(ctx, burrowID, false); err != nil {
		return nil, fmt.Errorf("failed to update burrow %d occupancy: %w", burrowID, err)
	}

	g.stats.AddAvailableBurrow()
	burrow.IsOccupied = false
	return burrow, nil
}

func (g *GopherApp) GetBurrowStatus(ctx context.Context) ([]*ent.Burrow, error) {
	return g.repo.GetAllBurrows(ctx)
}
