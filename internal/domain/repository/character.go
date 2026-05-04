package repository

import (
	"weakRaider/internal/domain/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CharacterRepository struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (c *CharacterRepository) Create(character *entity.Character) error {
	return c.db.Model(entity.Character{}).Create(character).Error
}

func (c *CharacterRepository) Update(character *entity.Character) error {
	return c.db.Model(&entity.Character{}).UpdateColumns(character).Error
}

func (c *CharacterRepository) Upsert(character *entity.Character) error {
	return c.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&character).Error
}

func (c *CharacterRepository) FindAll() (entity.CharacterMap, error) {
	characters := []entity.Character{}
	charMap := entity.CharacterMap{}

	if err := c.db.Model(entity.Character{}).Find(&characters).Error; err != nil {
		return nil, err
	}

	for _, v := range characters {
		charMap[v.Name+v.Realm] = v
	}

	return charMap, nil
}
