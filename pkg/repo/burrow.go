package repo

import "gophernet/pkg/db"

type IBurrowRepository interface {
}

type BurrowRepository struct {
	db db.Database
}

func NewBurrowRepository(db db.Database) *BurrowRepository {
	return &BurrowRepository{
		db: db,
	}
}
