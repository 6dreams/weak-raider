package main

import (
	"github.com/rs/zerolog/log"
	baseApp "weakRaider/internal/app"
)

func main() {
	app := baseApp.New()
	if err := app.Initialize(); err != nil {
		log.Fatal().Err(err).Msg("app initialization failed")

		return
	}

	log.Info().Msg("app initialized, starting.")

	app.Run()
}
