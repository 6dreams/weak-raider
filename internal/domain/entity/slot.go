package entity

type Slot struct {
	// Идентификатор слота, генерируется автоматически.
	ID int `gorm:"primary_key;column:id;type:bigint;primaryKey;autoIncrement;not null"`

	// Название слота экипировки.
	Name string `gorm:"column:name;type:text;not null"`
}

func (s Slot) TableName() string {
	return "slot"
}
