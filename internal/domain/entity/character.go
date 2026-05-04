package entity

import (
	"time"
)

type CharacterMap = map[string]Character

type Character struct {
	// Идентификатор персонажа, генерируется автоматически.
	ID int `gorm:"primary_key;column:id;type:bigserial;not null"`

	//Информация о персонаже, берётся из API WoWAudit.
	Name  string `gorm:"column:name;type:text;not null"`
	Class string `gorm:"column:class;type:text;not null"`
	Rank  string `gorm:"column:rank;type:text;not null"`
	Role  string `gorm:"column:role;type:text;not null"`
	Realm string `gorm:"column:realm;type:text;not null"`
	Note  string `gorm:"column:note;type:text;default:'';not null"`

	// Связь с гильдией.
	GuildId int
	Guild   *Guild `gorm:"foreignkey:guild_id"`

	// Время последнего обновления информации о персонаже.
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`

	// Последнее время обновления WoWAudit.
	// WowauditUpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;null"`
}

func (c Character) TableName() string {
	return "character"
}
