package entity

import "database/sql"

type Season struct {
	ID         int           `gorm:"primary_key;column:id;type:bigserial;not null"`
	BlizzardId sql.NullInt64 `gorm:"column:blizzard_id;type:bigint;null"`
	WowAuditId sql.NullInt64 `gorm:"column:woow_audit_id;type:bigint;null"`
	Name       string        `gorm:"column:name;type:text;not null"`
}

func (s Season) TableName() string {
	return "season"
}
