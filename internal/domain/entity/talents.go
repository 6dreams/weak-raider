package entity

import "weakRaider/internal/domain/entity/types"

type Talents struct {
	// Идентификатор набора талантов, генерируется автоматически.
	Id int `gorm:"primary_key;column:id;type:bigserial;not null"`

	// Название класса
	ClassName *types.Translation `gorm:"column:class_name;type:jsonb;not null"`
	ClassID   int                `gorm:"uniqueIndex:idx_player_talent;column:class_id;type:int;not null"`

	// Название специализации
	SpecName *types.Translation `gorm:"column:spec_name;type:jsonb;not null"`
	SpecID   int                `gorm:"uniqueIndex:idx_player_talent;column:spec_id;type:int;not null"`

	// Название героической ветки талантов
	HeroSpecName *types.Translation `gorm:"column:hero_spec_name;type:jsonb;not null"`
	HeroSpecID   int                `gorm:"uniqueIndex:idx_player_talent;column:hero_spec_id;type:int;not null"`

	// Общий набор талантов map[spellID]TalentNode, хранящийся в формате JsonB
	Talents map[int]TalentNode `gorm:"column:talents;type:jsonb;not null"`
}

type TalentNode struct {
	Ranks int                `json:"ranks"`
	Name  *types.Translation `json:"name"`
}

func (t Talents) TableName() string {
	return "talents"
}
