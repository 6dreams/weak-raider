package repository

import (
	"weakRaider/internal/domain/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TalentsRepository struct {
	db *gorm.DB
}

func NewTalentsRepository(db *gorm.DB) *TalentsRepository {
	return &TalentsRepository{db: db}
}

func (t *TalentsRepository) Upsert(talents *entity.Talents) error {
	return t.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(talents).Error
}

// GORM Hook
// func (t *TalentsRepository) AfterFind(tx *gorm.DB) error {

// 	return nil
// }
