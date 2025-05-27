package repo

import (
	"context"
	"fmt"
	"time"

	"gophernet/pkg/db"
	"gophernet/pkg/db/ent"
	"gophernet/pkg/db/ent/burrow"
)

// IBurrowRepository defines the interface for burrow data operations
type IBurrowRepository interface {
	GetAllBurrows(ctx context.Context) ([]*ent.Burrow, error)
	GetOccupiedBurrows(ctx context.Context) ([]*ent.Burrow, error)
	GetBurrowByID(ctx context.Context, id int) (*ent.Burrow, error)
	UpdateBurrowOccupancy(ctx context.Context, id int, isOccupied bool) error
	UpdateBurrowDepth(ctx context.Context, id int64, depth float64) error
	DeleteBurrow(ctx context.Context, id int64) error
	CreateBurrow(ctx context.Context, name string, depth float64, width float64, isOccupied bool, age int) (*ent.Burrow, error)
	CreateBurrows(ctx context.Context, burrows []*ent.Burrow) ([]*ent.Burrow, error)
	DeleteAllBurrows(ctx context.Context) error
}

// BurrowRepository implements the burrow data operations
type BurrowRepository struct {
	db db.Database
}

// NewBurrowRepository creates a new instance of BurrowRepository
func NewBurrowRepository(db db.Database) *BurrowRepository {
	return &BurrowRepository{
		db: db,
	}
}

// GetOccupiedBurrows retrieves all occupied burrows
func (r *BurrowRepository) GetOccupiedBurrows(ctx context.Context) ([]*ent.Burrow, error) {
	burrows, err := r.db.EntClient().Burrow.Query().
		Where(burrow.IsOccupied(true)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get occupied burrows: %w", err)
	}
	return burrows, nil
}

// UpdateBurrowDepth updates a burrow's depth and increments its age
func (r *BurrowRepository) UpdateBurrowDepth(ctx context.Context, id int64, depth float64) error {
	_, err := r.db.EntClient().Burrow.UpdateOneID(int(id)).
		SetDepth(depth).
		AddAge(1). // Increment age by 1 minute
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update burrow depth: %w", err)
	}
	return nil
}

// DeleteBurrow removes a burrow by ID
func (r *BurrowRepository) DeleteBurrow(ctx context.Context, id int64) error {
	if err := r.db.EntClient().Burrow.DeleteOneID(int(id)).Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete burrow: %w", err)
	}
	return nil
}

// CreateBurrow creates a new burrow
func (r *BurrowRepository) CreateBurrow(ctx context.Context, name string, depth float64, width float64, isOccupied bool, age int) (*ent.Burrow, error) {
	now := time.Now()
	burrow, err := r.db.EntClient().Burrow.Create().
		SetName(name).
		SetDepth(depth).
		SetWidth(width).
		SetIsOccupied(isOccupied).
		SetAge(age).
		SetUpdatedAt(now).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create burrow: %w", err)
	}
	return burrow, nil
}

// DeleteAllBurrows removes all burrows from the database
func (r *BurrowRepository) DeleteAllBurrows(ctx context.Context) error {
	if _, err := r.db.EntClient().Burrow.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete all burrows: %w", err)
	}
	return nil
}

// GetAllBurrows retrieves all burrows
func (r *BurrowRepository) GetAllBurrows(ctx context.Context) ([]*ent.Burrow, error) {
	burrows, err := r.db.EntClient().Burrow.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all burrows: %w", err)
	}
	return burrows, nil
}

// GetBurrowByID retrieves a burrow by its ID
func (r *BurrowRepository) GetBurrowByID(ctx context.Context, id int) (*ent.Burrow, error) {
	burrow, err := r.db.EntClient().Burrow.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get burrow by ID: %w", err)
	}
	return burrow, nil
}

// UpdateBurrowOccupancy updates a burrow's occupancy status
func (r *BurrowRepository) UpdateBurrowOccupancy(ctx context.Context, id int, isOccupied bool) error {
	_, err := r.db.EntClient().Burrow.UpdateOneID(id).
		SetIsOccupied(isOccupied).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update burrow occupancy: %w", err)
	}
	return nil
}

// CreateBurrows creates multiple burrows in a single transaction
func (r *BurrowRepository) CreateBurrows(ctx context.Context, burrows []*ent.Burrow) ([]*ent.Burrow, error) {
	now := time.Now()
	bulk := make([]*ent.BurrowCreate, len(burrows))
	for i, b := range burrows {
		bulk[i] = r.db.EntClient().Burrow.Create().
			SetName(b.Name).
			SetDepth(b.Depth).
			SetWidth(b.Width).
			SetIsOccupied(b.IsOccupied).
			SetAge(b.Age).
			SetUpdatedAt(now)
	}
	createdBurrows, err := r.db.EntClient().Burrow.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create burrows in bulk: %w", err)
	}
	return createdBurrows, nil
}
