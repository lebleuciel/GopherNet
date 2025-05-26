package app

import (
	"context"
	"errors"
	"fmt"

	"gophernet/pkg/db/ent"
	"gophernet/pkg/repo"
)

type IGopherApp interface {
	GetGopher(ctx context.Context) (string, error)
	StartScheduler(ctx context.Context)
	StopScheduler()
	RentBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error)
	ReleaseBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error)
	GetBurrowStatus(ctx context.Context) ([]*ent.Burrow, error)
}

type GopherApp struct {
	repo      repo.IBurrowRepository
	scheduler *Scheduler
}

func NewGopherApp(repo repo.IBurrowRepository) *GopherApp {
	return &GopherApp{
		repo:      repo,
		scheduler: NewScheduler(repo),
	}
}

func (g *GopherApp) GetGopher(ctx context.Context) (string, error) {
	return "Gopher is ready to help!", nil
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
		return nil, fmt.Errorf("failed to get burrow: %w", err)
	}

	if burrow.IsOccupied {
		return nil, errors.New("burrow is already occupied")
	}

	if err := g.repo.UpdateBurrowOccupancy(ctx, burrowID, true); err != nil {
		return nil, fmt.Errorf("failed to update burrow occupancy: %w", err)
	}

	burrow.IsOccupied = true
	return burrow, nil
}

func (g *GopherApp) ReleaseBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error) {
	burrow, err := g.repo.GetBurrowByID(ctx, burrowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get burrow: %w", err)
	}

	if !burrow.IsOccupied {
		return nil, errors.New("burrow is not occupied")
	}

	if err := g.repo.UpdateBurrowOccupancy(ctx, burrowID, false); err != nil {
		return nil, fmt.Errorf("failed to update burrow occupancy: %w", err)
	}

	burrow.IsOccupied = false
	return burrow, nil
}

func (g *GopherApp) GetBurrowStatus(ctx context.Context) ([]*ent.Burrow, error) {
	burrows, err := g.repo.GetAllBurrows(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get burrows: %w", err)
	}
	return burrows, nil
}
