package main

import (
	"dothething/internal/client"
	"dothething/internal/cmd"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// logger configuration
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// initializing the client
	clientAPI, err := client.NewAPIClient()
	if err != nil {
		log.Error().Err(err)
	}

	err = cmd.New(clientAPI).Run()
	if err != nil {
		log.Error().Err(err)
	}
}
