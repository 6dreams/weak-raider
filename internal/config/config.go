package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

// Database - contains all parameters database connection.
type Database struct {
	Host        string `env:"WR_PG_HOST"`
	Port        string `env:"WR_PG_PORT"`
	User        string `env:"WR_PG_USER"`
	Password    string `env:"WR_PG_PASS"`
	Migrations  string `yaml:"migrations"`
	Name        string `env:"WR_PG_NAME"`
	SslMode     string `yaml:"sslmode"`
	Driver      string `yaml:"driver"`
	Connections DBCons `yaml:"connections"`
}

type DBCons struct {
	MaxOpenCons     int           `yaml:"maxOpenCons"`
	MaxIdleCons     int           `yaml:"maxIdleCons"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Port        string `yaml:"port"`
	Version     string
	CommitHash  string
}

type Secrets struct {
	Blizzard struct {
		Client string `env:"WR_BLIZZARD_CLIENT"`
		Secret string `env:"WR_BLIZZARD_SECRET"`
	}
	Logs struct {
		Client string `env:"WR_WARCRAFTLOGS_CLIENT"`
		Secret string `env:"WR_WARCRAFTLOGS_SECRET"`
	}
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project  Project  `yaml:"project"`
	Database Database `yaml:"database"`
	Auth     Secrets
}

// ReadConfig - read configurations from file and init instance Config.
func ReadConfig(filePath string) (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	var cfg Config
	if err := cleanenv.ReadConfig(filePath, &cfg); err != nil {
		return nil, err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return &cfg, nil
}
