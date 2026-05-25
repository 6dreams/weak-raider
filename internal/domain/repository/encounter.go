package repository

import (
	"weakRaider/internal/domain/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EncounterRepository struct {
	db *gorm.DB
}

func NewEncounterRepository(db *gorm.DB) *EncounterRepository {
	return &EncounterRepository{db: db}
}

func (r *EncounterRepository) Upsert(encounter *entity.Encounter) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&encounter).Error
}

func (r *EncounterRepository) Update(encounter *entity.Encounter) error {
	return r.db.Save(encounter).Error
}
