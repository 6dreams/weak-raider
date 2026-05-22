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
	migration := flag.Bool("migration", true, "Defines the migration start option")

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}

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
			log.Error().Err(err).Msg("Migration-up failed")
			return
		}
	} else {
		if err = goose.Down(sqlDB, cfg.Database.Migrations); err != nil {
			log.Error().Err(err).Msg("Migration-down failed")
			return
		}
	}
}
