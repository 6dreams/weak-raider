package entity

import "database/sql"

type SeasonMap = map[int]Season

type Season struct {
	// Идентификатор сезона, генерируется автоматически.
	ID int `gorm:"primary_key;column:id;type:bigserial;not null"`

	// Идентификатор сезона по Blizzard, берётся из API Blizzard.
	BlizzardId int `gorm:"column:blizzard_id;type:bigint;not null"`

	// Идентификатор сезона по WoWAudit, берётся из API WoWAudit.
	WowAuditId sql.NullInt64 `gorm:"column:wow_audit_id;type:bigint;null"`

	// Название сезона (берётся из WoWAudit)
	Name string `gorm:"column:name;type:text;not null"`
}

func (s Season) TableName() string {
	return "season"
}
