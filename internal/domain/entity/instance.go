package entity

import (
	"time"
	"weakRaider/internal/domain/entity/types"
)

type InstanceMap = map[int]Instance

type Instance struct {
	// Идентификатор подземелья, берётся из API Blizzard.
	ID int `gorm:"primary_key;column:id;type:bigint;primaryKey;autoIncrement;not null"`

	// Связь с сезоном.
	SeasonId int
	Season   *Season `gorm:"foreignkey:season_id"`

	Encounters []Encounter `gorm:"references:Id"`

	// Название подземелья.
	Name *types.Translation `gorm:"column:name;type:jsonb;not null"`

	// Подземелье является рейдом.
	IsRaid bool `gorm:"column:is_raid;type:boolean;not null"`

	// Время последнего обновления.
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (i Instance) TableName() string {
	return "instance"
}

func (i Instance) GetEncounter(id int) *Encounter {
	for _, encounter := range i.Encounters {
		if encounter.ID == id {
			return &encounter
		}
	}

	return nil
}
