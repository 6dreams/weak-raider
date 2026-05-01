package repo

import (
	"context"
	"weakRaider/internal/model"

	"gorm.io/gorm"
)

type Repo interface {
	Add(ctx context.Context, char model.Character) error
	Upsert(ctx context.Context, char model.Character) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}
