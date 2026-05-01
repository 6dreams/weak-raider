package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
	"weakRaider/internal/config"
	"weakRaider/internal/domain/repository"
	appLogger "weakRaider/internal/logger"
)

type App struct {
	Config *config.Config
	Logger *zerolog.Logger

	db *gorm.DB

	Repository struct {
		Season *repository.SeasonRepository
	}
}

func New() *App {
	return &App{}
}

func (app *App) Run() {
	// do job.
}

func (app *App) Database() *gorm.DB {
	return app.db
}

func (app *App) Initialize() error {
	// config
	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("config load error: %v", err)
	}
	app.Config = cfg

	// logger
	logger := appLogger.LogInit(cfg.Project.Debug)
	app.Logger = &logger

	// database
	app.db, err = app.configureGorm()
	if err != nil {
		return fmt.Errorf("configure Gorm error: %v", err)
	}

	return nil
}

func (app *App) configureGorm() (*gorm.DB, error) {
	gormConfig := &gorm.Config{}

	conn := postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
			app.Config.Database.Host,
			app.Config.Database.Port,
			app.Config.Database.User,
			app.Config.Database.Password,
			app.Config.Database.Name,
			app.Config.Database.SslMode,
		),
	})

	// sql logger
	if app.Config.Project.Debug {
		gormConfig.Logger = gormLogger.New(
			app.Logger,
			gormLogger.Config{
				SlowThreshold:             20 * time.Millisecond,
				LogLevel:                  gormLogger.Info,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
			},
		)
	}

	db, err := gorm.Open(conn, gormConfig)
	if err != nil {
		return nil, err
	}

	// repositories
	app.Repository.Season = repository.NewSeasonRepository(app.db)

	return db, nil
}
