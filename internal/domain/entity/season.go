package entity

import (
	"database/sql"
	"time"
)

type SeasonMap = map[int]Season

type Season struct {
	// Идентификатор сезона, генерируется автоматически.
	ID int `gorm:"primary_key;column:id;type:bigint;primaryKey;autoIncrement;not null"`

	// Идентификатор сезона по Blizzard, берётся из API Blizzard.
	BlizzardId int `gorm:"column:blizzard_id;type:bigint;not null"`

	// Идентификатор сезона по WoWAudit, берётся из API WoWAudit.
	WowAuditId sql.NullInt64 `gorm:"column:wow_audit_id;type:bigint;null"`

	// Название сезона (берётся из WoWAudit)
	Name string `gorm:"column:name;type:text;not null"`

	// Сезон является текущим.
	IsCurrent bool `gorm:"column:is_current;type:boolean;not null"`

	// Время последнего обновления.
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (s Season) TableName() string {
	return "season"
}
