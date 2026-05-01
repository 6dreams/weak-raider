package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"weakRaider/internal/domain/entity"
)

type SeasonRepository struct {
	db *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) *SeasonRepository {
	return &SeasonRepository{db: db}
}

func (r *SeasonRepository) Upsert(season *entity.Season) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(season).Error
}
