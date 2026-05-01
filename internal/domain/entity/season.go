package entity

import "database/sql"

type Season struct {
	// Идентификатор сезона, генерируется автоматически.
	ID int `gorm:"primary_key;column:id;type:bigserial;not null"`

	// Идентификатор сезона по Blizzard, берётся из API Blizzard.
	BlizzardId sql.NullInt64 `gorm:"column:blizzard_id;type:bigint;null"`

	// Идентификатор сезона по WoW Audit, берётся из API WoW Audit.
	WowAuditId sql.NullInt64 `gorm:"column:wow_audit_id;type:bigint;null"`

	// Название сезона (берётся из wowaudit)
	Name string `gorm:"column:name;type:text;not null"`
}

func (s Season) TableName() string {
	return "season"
}
