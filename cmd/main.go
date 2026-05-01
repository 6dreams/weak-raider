package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"weakRaider/internal/app"
)

var (
	logger zerolog.Logger
)

func main() {
	application := app.New()
	if err := application.Initialize(); err != nil {
		log.Fatal().Err(err).Msg("application initialization failed")

		return
	}
	application.Run()

	//cfg, err := config.ReadConfig()
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Failed init configuration")
	//}
	//
	//logger = myLogger.LogInit(cfg.Project.Debug)
	//
	//dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
	//	cfg.Database.Host,
	//	cfg.Database.Port,
	//	cfg.Database.User,
	//	cfg.Database.Password,
	//	cfg.Database.Name,
	//	cfg.Database.SslMode,
	//)
	//
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//logger.Debug().Msg("gorm connection succesfully created")
	//if err != nil {
	//	fmt.Println("-")
	//	logger.Fatal().Err(err).Msg("failed init gorm")
	//}
	//
	//sqlDB, err := db.DB()
	//logger.Debug().Msg("gorm connection succesfully created")
	//
	//if err != nil {
	//	logger.Fatal().Err(err).Msg("failed get sqlDB from gorm")
	//}
	//defer sqlDB.Close()
	//
	//repo := repo.NewRepo(db)
	//service := service.NewService(repo, logger)
	//handlers := handlers.NewHandler(service)
	//server := server.NewServer(cfg.Project.Port, handlers.Handler())
	//
	//server.Run()
}
