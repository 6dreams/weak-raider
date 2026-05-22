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

func (t *TalentsRepository) Create(TalentsJsonB *entity.TalentsJsonB) error {
	return t.db.Model(entity.TalentsJsonB{}).Create(TalentsJsonB).Error
}

func (t *TalentsRepository) Update(TalentsJsonB *entity.TalentsJsonB) error {
	return t.db.Model(&entity.TalentsJsonB{}).UpdateColumns(TalentsJsonB).Error
}

// func (t *TalentsRepository) Find(class, spec, heroSpec string) error {
// 	return t.db.Model(&entity.TalentsJsonB{}).Find()
// }

func (t *TalentsRepository) Upsert(TalentsJsonB *entity.TalentsJsonB) error {
	return t.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&TalentsJsonB).Error
}

// GORM Hook
func (t *TalentsRepository) AfterFind(tx *gorm.DB) error {

	return nil
}
