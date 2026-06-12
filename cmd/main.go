package main

import (
	baseApp "weakRaider/internal/app"

	"github.com/rs/zerolog/log"
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
