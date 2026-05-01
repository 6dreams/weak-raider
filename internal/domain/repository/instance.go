package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"weakRaider/internal/domain/entity"
)

type InstanceRepository struct {
	db *gorm.DB
}

func NewInstanceRepository(db *gorm.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) FindAll() (entity.InstanceMap, error) {
	instances := entity.InstanceMap{}
	if err := r.db.Model(entity.Instance{}).Find(&instances).Error; err != nil {
		return nil, err
	}

	return instances, nil
}

func (r *InstanceRepository) Upsert(instance *entity.Instance) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&instance).Error
}
