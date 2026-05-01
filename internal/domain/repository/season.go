package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"weakRaider/internal/domain/entity"
)

var ErrSeasonsNotFound = errors.New("seasons not found")

type SeasonRepository struct {
	db *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) *SeasonRepository {
	return &SeasonRepository{db: db}
}

func (r *SeasonRepository) Upsert(season *entity.Season) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&season).Error
}

func (r *SeasonRepository) FindAll() (entity.SeasonMap, error) {
	var seasons []entity.Season
	if err := r.db.Find(&seasons).Error; err != nil {
		return nil, ErrSeasonsNotFound
	}

	mapped := entity.SeasonMap{}
	for _, season := range seasons {
		mapped[season.BlizzardId] = season
	}

	return mapped, nil
}
