package entity

import "time"

type Guild struct {
	// Идентификатор гильдии, генерируется автоматически.
	ID int `gorm:"primary_key;column:id;type:bigserial;not null"`

	// Название гильдии.
	Name string `gorm:"column:name;type:text;not null"`

	// Приватный ключ для работы с API WoW Audit.
	WowauditKey string `gorm:"column:wow_audit_key;type:text;not null"`

	//Время последнего обновления информации о гильдии.
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (g Guild) TableName() string {
	return "guild"
}
