package entity

import (
	"database/sql"
	"time"
)

type Encounter struct {
	// Идентификатор боя, генерируется автоматическси.
	ID int `gorm:"primary_key;column:id;type:bigint;primaryKey;autoIncrement;not null"`

	// Имя босса.
	Name string `gorm:"column:name;type:text;not null"`

	BlizzardId int `gorm:"column:blizzard_id;type:int;not null"`

	InstanceId int64
	Instance   *Instance `gorm:"foreignkey:instance_id"`

	// Идентификатор боя, берётся из API WarcraftLogs.
	LogsEncounterID sql.NullInt64 `gorm:"column:logs_encounter_id;type:bigint;null"`

	// Время последнего обновления информации о бое.
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (e Encounter) TableName() string {
	return "encounter"
}
