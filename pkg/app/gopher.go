package app

import (
	"context"

	"gophernet/pkg/db/ent"
	apperrors "gophernet/pkg/errors"
	"gophernet/pkg/logger"
	"gophernet/pkg/repo"

	"go.uber.org/zap"
)

type IGopherApp interface {
	GetGopher(ctx context.Context) (string, error)
	RentBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error)
	ReleaseBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error)
	GetBurrowStatus(ctx context.Context) ([]*ent.Burrow, error)
	GetBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error)
}

type GopherApp struct {
	repo repo.IBurrowRepository
	log  *zap.Logger
}

func NewGopherApp(repo repo.IBurrowRepository) *GopherApp {
	ga := &GopherApp{
		repo: repo,
		log:  logger.Get(),
	}
	return ga
}

func (g *GopherApp) GetGopher(ctx context.Context) (string, error) {
	g.log.Debug("Getting gopher status")
	return "Gopher is ready to help!", nil
}

func (g *GopherApp) RentBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error) {
	g.log.Info("Attempting to rent burrow", zap.Int("burrow_id", burrowID))

	burrow, err := g.repo.GetBurrowByID(ctx, burrowID)
	if err != nil {
		g.log.Error("Failed to get burrow", zap.Int("burrow_id", burrowID), zap.Error(err))
		return nil, err
	}

	if burrow.IsOccupied {
		g.log.Warn("Burrow is already occupied", zap.Int("burrow_id", burrowID))
		return nil, apperrors.ErrBurrowOccupied
	}

	if err := g.repo.UpdateBurrowOccupancy(ctx, burrowID, true); err != nil {
		g.log.Error("Failed to update burrow occupancy", zap.Int("burrow_id", burrowID), zap.Error(err))
		return nil, err
	}

	burrow.IsOccupied = true
	g.log.Info("Successfully rented burrow", zap.Int("burrow_id", burrowID))
	return burrow, nil
}

func (g *GopherApp) ReleaseBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error) {
	g.log.Info("Attempting to release burrow", zap.Int("burrow_id", burrowID))

	burrow, err := g.repo.GetBurrowByID(ctx, burrowID)
	if err != nil {
		g.log.Error("Failed to get burrow", zap.Int("burrow_id", burrowID), zap.Error(err))
		return nil, err
	}

	if !burrow.IsOccupied {
		g.log.Warn("Burrow is not occupied", zap.Int("burrow_id", burrowID))
		return nil, apperrors.ErrBurrowNotOccupied
	}

	if err := g.repo.UpdateBurrowOccupancy(ctx, burrowID, false); err != nil {
		g.log.Error("Failed to update burrow occupancy", zap.Int("burrow_id", burrowID), zap.Error(err))
		return nil, err
	}

	burrow.IsOccupied = false
	g.log.Info("Successfully released burrow", zap.Int("burrow_id", burrowID))
	return burrow, nil
}

func (g *GopherApp) GetBurrowStatus(ctx context.Context) ([]*ent.Burrow, error) {
	g.log.Debug("Getting burrow status")

	burrows, err := g.repo.GetAllBurrows(ctx)
	if err != nil {
		g.log.Error("Failed to get burrows", zap.Error(err))
		return nil, apperrors.Wrap(err, "failed to get burrows")
	}

	g.log.Info("Retrieved burrow status", zap.Int("count", len(burrows)))
	return burrows, nil
}

func (g *GopherApp) GetBurrow(ctx context.Context, burrowID int) (*ent.Burrow, error) {
	g.log.Debug("Getting burrow", zap.Int("burrow_id", burrowID))

	burrow, err := g.repo.GetBurrowByID(ctx, burrowID)
	if err != nil {
		g.log.Error("Failed to get burrow", zap.Int("burrow_id", burrowID), zap.Error(err))
		return nil, err
	}

	g.log.Info("Retrieved burrow", zap.Int("burrow_id", burrowID))
	return burrow, nil
}
