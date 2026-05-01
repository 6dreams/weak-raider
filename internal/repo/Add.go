package repo

import (
	"context"
	"weakRaider/internal/model"
)

func (r *repo) Add(ctx context.Context, char model.Character) error {
	tx := r.db.Create(&char)
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Commit()
	if result.Error != nil {
		return result.Error
	}
	return nil
}
