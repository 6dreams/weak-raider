package main

import (
	"fmt"
	"log"
	"strings"
	baseApp "weakRaider/internal/app"
	"weakRaider/internal/domain/entity"

	"gorm.io/gorm"
)

var entities = []interface{}{
	entity.Character{},
	entity.Encounter{},
	entity.Guild{},
	entity.Instance{},
	entity.Item{},
	entity.Season{},
	entity.Slot{},
	entity.TalentsJsonB{},
}

func main() {
	app := baseApp.New()
	if err := app.Initialize(); err != nil {
		log.Panicf("failed to initialize app: %v", err)
	}

	db := app.Database()

	tx := db.Begin()

	var stmts []string
	err := tx.Callback().Raw().Register("schema diffs", func(tx *gorm.DB) {
		stmts = append(stmts, tx.Statement.SQL.String())
	})

	if err != nil {
		log.Panicf("failed to register schema diff callback: %v", err)
	}

	if err := tx.AutoMigrate(entities...); err != nil {
		log.Panicf("failed to auto migrate: %v", err)
	}

	tx.Rollback()

	if len(stmts) > 0 {
		fmt.Println("\nDatabase schema diff:")
		fmt.Println(strings.Join(stmts, "\n"))
	} else {
		fmt.Println("\nDatabase schema has no changes")
	}
}
