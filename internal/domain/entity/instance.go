package entity

import "time"

type Instance struct {
	// Идентификатор подземелья, берётся из API Blizzard.
	ID int `gorm:"primary_key;column:id;type:bigint;not null"`

	// Связь с сезоном.
	SeasonId int
	Season   *Season `gorm:"foreignkey:season_id"`

	// Название подземелья.
	Name string `gorm:"column:name;type:text;not null"`

	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (i Instance) TableName() string {
	return "instance"
}
