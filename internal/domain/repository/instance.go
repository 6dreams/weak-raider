package repository

import (
	"weakRaider/internal/domain/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InstanceRepository struct {
	db *gorm.DB
}

func NewInstanceRepository(db *gorm.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) FindAll() (entity.InstanceMap, error) {
	instances := make(entity.InstanceMap)
	var instanceList []entity.Instance
	if err := r.db.Model(entity.Instance{}).Find(&instanceList).Error; err != nil {
		return nil, err
	}

	for _, instance := range instanceList {
		instances[instance.ID] = instance
	}

	return instances, nil
}

func (r *InstanceRepository) FindWithInstances() ([]entity.Instance, error) {
	instances := make([]entity.Instance, 0)
	if err := r.db.Model(entity.Instance{}).Preload("Encounters").Find(&instances).Error; err != nil {
		return nil, err
	}

	return instances, nil
}

func (r *InstanceRepository) Upsert(instance *entity.Instance) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&instance).Error
}

func (r *InstanceRepository) Update(instance *entity.Instance) error {
	return r.db.Save(&instance).Error
}
