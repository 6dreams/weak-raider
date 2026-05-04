package config

import (
	"flag"
	"fmt"
	"os"
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
	Debug       bool          `yaml:"debug"`
	Name        string        `yaml:"name"`
	Environment string        `yaml:"environment"`
	Port        string        `yaml:"port"`
	Tickrate    time.Duration `yaml:"tickrate"`
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
func ReadConfig() (*Config, error) {
	err, env, config := getConfigFiles()
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(env); err != nil {
		return nil, err
	}

	var cfg Config
	if err := cleanenv.ReadConfig(config, &cfg); err != nil {
		return nil, err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return &cfg, nil
}

func getConfigFiles() (error, string, string) {
	var env, conf string

	flag.StringVar(&env, "env", ".env", "path to environment file")
	flag.StringVar(&conf, "config", "", "path to config file")
	flag.Parse()

	if conf == "" {
		return fmt.Errorf("config file not set. Use -config"), "", ""
	}

	if _, err := os.Stat(env); os.IsNotExist(err) {
		return fmt.Errorf("environment file does not exist: %v", conf), "", ""
	}

	if _, err := os.Stat(conf); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %v", conf), "", ""
	}

	return nil, env, conf
}
