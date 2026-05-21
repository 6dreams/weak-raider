package repository

import "gorm.io/gorm"

type EncounterRepository struct {
	db *gorm.DB
}
