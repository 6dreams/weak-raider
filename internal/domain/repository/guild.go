package repository

import (
	"weakRaider/internal/domain/entity"

	"gorm.io/gorm"
)

type GuildRepository struct {
	db *gorm.DB
}

func NewGuildRepository(db *gorm.DB) *GuildRepository {
	return &GuildRepository{db: db}
}

func (c *GuildRepository) Create(guild *entity.Guild) error {
	return c.db.Model(entity.Guild{}).Create(guild).Error
}

func (c *GuildRepository) FindAll() ([]entity.Guild, error) {
	guilds := []entity.Guild{}

	if err := c.db.Model(entity.Guild{}).Find(&guilds).Error; err != nil {
		return nil, err
	}

	return guilds, nil
}
