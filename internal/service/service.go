package service

import (
	"context"
	"weakRaider/internal/model"

	"github.com/rs/zerolog"
)

type Repo interface {
	Add(ctx context.Context, char model.Character) error
	Upsert(ctx context.Context, char model.Character) error
}

type service struct {
	repo   Repo
	logger zerolog.Logger
}

func NewService(repo Repo, logger zerolog.Logger) *service {
	return &service{repo: repo, logger: logger}
}

// func (s *service) Add(ctx context.Context, char model.Character) error {
// 	return s.repo.Add(ctx, char)
// }

// func (s *service) Upsert(ctx context.Context, char model.Character) error {
// 	return s.repo.Upsert(ctx, char)
// }
