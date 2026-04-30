package myLogger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogInit(debugLevel bool) zerolog.Logger {
	if debugLevel {
		return log.Level(zerolog.DebugLevel)
	}
	return log.Level(zerolog.InfoLevel)
}
