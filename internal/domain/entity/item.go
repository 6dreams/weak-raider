package entity

type Item struct {
	// Идентификатор предмета, берётся из API Blizzard.
	ID int `gorm:"primary_key;column:id;type:bigint;primaryKey;autoIncrement;not null"` // ign id

	// Характеристики предмета.

	// Название предмета.
	Name string `gorm:"column:name;type:text;not null"`

	// Является ли уникальным.
	IsUnique bool `gorm:"column:is_unique;type:bool;not null"`

	// Количество сокетов.
	Sockets int `gorm:"column:sockets;type:int;null"`

	// Связь с подземельем.
	InstanceId int
	Instance   *Instance `gorm:"foreignkey:instance_id"`

	// связь с боссом.
	EncounterId int
	Encounter   *Encounter `gorm:"foreignkey:encounter_id"`

	// связь со слотом экипировки.
	SlotId int
	Slot   *Slot `gorm:"foreignkey:slot_id"`
}

func (i Item) TableName() string {
	return "item"
}
