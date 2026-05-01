package entity

import "time"

type Encounter struct {
	// Идентификатор боя, генерируется автоматическси.
	ID int `gorm:"primary_key;column:id;type:bigserial;not null"`

	// Имя босса.
	Name string `gorm:"column:name;type:text;not null"`

	// Идентификатор боя, берётся из API Warcraft Logs.
	LogsEncounterID int `gorm:"column:logs_encounter_id;type:int;not null"`

	// Время последнего обновления информации о бое.
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (e Encounter) TableName() string {
	return "encounter"
}
