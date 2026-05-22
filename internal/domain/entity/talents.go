package entity

type TalentsJsonB struct {
	// Идентификатор набора талантов, генерируется автоматически.
	ID int `gorm:"primary_key;column:id;type:bigserial;not null"`

	// Название класса
	ClassName string `gorm:"column:class_name;type:text;not null"`

	// Название специализации
	SpecName string `gorm:"column:specialisation_name;type:text;not null"`

	// Название героической ветки талантов
	HeroSpecName string `gorm:"column:hero_specialisation_name;type:text;not null"`

	// Общий набор талантов, хранящийся в формате JsonB
	talents []byte `gorm:"column:talents;type:jsonb;not null"`
}

type TalentsMap struct {
	ClassName    string
	SpecName     string
	HeroSpecName string
	Talents      map[int]string
}

func (t TalentsJsonB) TableName() string {
	return "talents"
}
