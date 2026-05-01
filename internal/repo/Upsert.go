package repo

import (
	"context"
	"weakRaider/internal/model"

	"gorm.io/gorm/clause"
)

func (r *repo) Upsert(ctx context.Context, char model.Character) error {
	tx := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"class",
			"rank",
			"role",
			"realm",
			"note",
			"guild_id",
			"updated_at",
		}),
	}).Create(&char)

	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Commit()
	if result.Error != nil {
		return result.Error
	}

	return nil
}
