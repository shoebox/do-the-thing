package main

import (
	"dothething/internal/api"
	"dothething/internal/client"
	"dothething/internal/cmd"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfg *api.Config

func main() {
	// logger configuration
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// intializing the client
	cfg = &api.Config{}
	clientAPI, err := client.NewAPIClient(cfg)
	if err != nil {
		log.Error().Err(err)
	}

	err = cmd.New(clientAPI, cfg).Run()
	if err != nil {
		log.Error().Err(err)
	}
}
