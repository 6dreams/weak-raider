package app

import (
	"fmt"
	"net/http"
	"time"
	"weakRaider/internal/clients"
	"weakRaider/internal/clients/blizzard"
	"weakRaider/internal/clients/warcraftlogs"
	"weakRaider/internal/clients/wowaudit"
	"weakRaider/internal/config"
	"weakRaider/internal/domain/manager"
	"weakRaider/internal/domain/repository"
	appLogger "weakRaider/internal/logger"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type App struct {
	Config *config.Config
	Logger *zerolog.Logger

	db *gorm.DB

	Repository struct {
		Season    *repository.SeasonRepository
		Instance  *repository.InstanceRepository
		Character *repository.CharacterRepository
	}

	Client struct {
		Blizzard  *blizzard.ClientWithResponses
		BattleNet *blizzard.ClientWithResponses
		WowAudit  *wowaudit.ClientWithResponses
		Logs      *warcraftlogs.ClientWithResponses
	}

	Manager struct {
		Auth          *manager.AuthManager
		SeasonSync    *manager.SeasonSync
		CharacterSync *manager.CharacterSync
	}

	Keys *clients.ApiKeys
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

	if err = app.configureClients(); err != nil {
		return fmt.Errorf("configure Clients error: %v", err)
	}

	app.configureManagers()

	return nil
}

func (app *App) configureClients() error {
	hc := http.Client{}
	var err error

	if app.Client.Blizzard, err = blizzard.NewClientWithResponses(blizzard.ServerUrlHttpseuApiBlizzardCom, blizzard.WithHTTPClient(&hc)); err != nil {
		return err
	}

	if app.Client.BattleNet, err = blizzard.NewClientWithResponses(blizzard.ServerUrlHttpsoauthBattleNet, blizzard.WithHTTPClient(&hc)); err != nil {
		return err
	}

	if app.Client.WowAudit, err = wowaudit.NewClientWithResponses(wowaudit.ServerUrlHttpswowauditCom, wowaudit.WithHTTPClient(&hc)); err != nil {
		return err
	}

	if app.Client.Logs, err = warcraftlogs.NewClientWithResponses(warcraftlogs.ServerUrlHttpswwwWarcraftlogsCom, warcraftlogs.WithHTTPClient(&hc)); err != nil {
		return err
	}

	return nil
}

func (app *App) configureManagers() {
	app.Manager.Auth = manager.NewAuthManager(
		app.Config,
		app.Client.BattleNet,
		app.Client.Logs,
		app.Keys,
	)
	app.Manager.SeasonSync = manager.NewSeasonSync(
		app.Repository.Season,
		app.Repository.Instance,
		app.Client.Blizzard,
		app.Manager.Auth,
	)
	app.Manager.CharacterSync = manager.NewCharacterSync(
		app.Repository.Character,
		app.Client.WowAudit,
	)
}

func (app *App) configureGorm() (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		DisableAutomaticPing: true,
	}

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
	app.Repository.Season = repository.NewSeasonRepository(db)
	app.Repository.Instance = repository.NewInstanceRepository(db)
	app.Repository.Character = repository.NewCharacterRepository(db)

	return db, nil
}
