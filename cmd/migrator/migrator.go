package main

import (
	"flag"
	"fmt"
	"weakRaider/internal/config"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}

	cfg := config.GetConfigInstance()
	migration := flag.Bool("migration", true, "Defines the migration start option")

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Debug().Msg("gorm connection succesfully created")
	if err != nil {
		fmt.Println("-")
		log.Fatal().Err(err).Msg("Failed init gorm")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed get sqlDB from gorm")
	}
	defer sqlDB.Close()

	if *migration {
		if err = goose.Up(sqlDB, cfg.Database.Migrations); err != nil {
			log.Error().Err(err).Msg("Migration failed")
			return
		}
	}
}
