package main

import (
	"fmt"
	"weakRaider/internal/config"
	myLogger "weakRaider/internal/logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	logger zerolog.Logger
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}

	cfg := config.GetConfigInstance()

	logger = myLogger.LogInit(cfg.Project.Debug)

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	logger.Debug().Msg("gorm connection succesfully created")
	if err != nil {
		fmt.Println("-")
		logger.Fatal().Err(err).Msg("failed init gorm")
	}

	sqlDB, err := db.DB()
	logger.Debug().Msg("gorm connection succesfully created")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed get sqlDB from gorm")
	}
	defer sqlDB.Close()

}
