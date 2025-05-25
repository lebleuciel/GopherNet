package repo

import (
	"context"
	"gophernet/pkg/db"
	"gophernet/pkg/db/ent"
	"gophernet/pkg/db/ent/burrow"
)

type IBurrowRepository interface {
	GetOccupiedBurrows(ctx context.Context) ([]*ent.Burrow, error)
	UpdateBurrowDepth(ctx context.Context, id int64, depth float64) error
	DeleteBurrow(ctx context.Context, id int64) error
	CreateBurrow(ctx context.Context, name string, depth float64, width float64, isOccupied bool, age int) (*ent.Burrow, error)
}

type BurrowRepository struct {
	db db.Database
}

func NewBurrowRepository(db db.Database) *BurrowRepository {
	return &BurrowRepository{
		db: db,
	}
}

func (r *BurrowRepository) GetOccupiedBurrows(ctx context.Context) ([]*ent.Burrow, error) {
	return r.db.EntClient().Burrow.Query().
		Where(burrow.IsOccupied(true)).
		All(ctx)
}

// UpdateBurrowDepth updates a burrow's depth and increments its age
func (r *BurrowRepository) UpdateBurrowDepth(ctx context.Context, id int64, depth float64) error {
	_, err := r.db.EntClient().Burrow.UpdateOneID(int(id)).
		SetDepth(depth).
		AddAge(1). // Increment age by 1 minute
		Save(ctx)
	return err
}

// DeleteBurrow deletes a burrow by ID
func (r *BurrowRepository) DeleteBurrow(ctx context.Context, id int64) error {
	return r.db.EntClient().Burrow.DeleteOneID(int(id)).Exec(ctx)
}

func (r *BurrowRepository) CreateBurrow(ctx context.Context, name string, depth float64, width float64, isOccupied bool, age int) (*ent.Burrow, error) {
	return r.db.EntClient().Burrow.Create().
		SetName(name).
		SetDepth(depth).
		SetWidth(width).
		SetIsOccupied(isOccupied).
		SetAge(age).
		Save(ctx)
}
