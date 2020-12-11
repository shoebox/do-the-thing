package main

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/client"
	"dothething/internal/cmd"
	"dothething/internal/xcode/project"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfg *api.Config

// var executor api.Executor
var targetProject project.Project

var clientAPI api.API

func main() {

	// logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

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

func c() {

	/*
		// Context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel() // The cancel should be deferred so resources are cleaned up

		cfg = api.Config{
			Path:          "",
			Scheme:        "",
			Configuration: "Demo Debug",
			Target:        "Axis_iOS",
			CodeSignOption: api.SignConfig{
				CertificatePassword: "abc12345",
				Path:                "/users/johann.martinache/Desktop/dummy",
			},
		}

		clientAPI, err := client.NewAPIClient(&cfg)
		if err != nil {
			panic(err)
		}

		pj, err := clientAPI.XCodeProjectService().Parse(ctx)
		for _, tgt := range pj.Pbx.Targets {
			fmt.Printf(">>>%#v\n", tgt.Dependencies)
		}

		// pj.ValidateConfiguration(cfg)

		// err = resolveSignature(ctx, pj)
		// fmt.Println("err", err)

		// keychainTest()
		// build()
		// archive()
		unitTest()
	*/
}

func listXCodeInstance(ctx context.Context) error {
	installs, err := clientAPI.XCodeListService().List(ctx)
	if err != nil {
		return err
	}

	for _, i := range installs {
		fmt.Println(i.Path, i.Version)
	}

	return nil
}

func selectService() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	res, err := clientAPI.XCodeSelectService().Find(ctx, "11.3.1")
	if err != nil {
		return err
	}

	fmt.Println("XCode instance : ", res, err)
	return nil
}

/*
func unitTest() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Retrieving the destination for the scheme
	dd, err := clientAPI.
		DestinationService().
		List(ctx, cfg.Scheme)

	if err != nil {
		return err
	}

	// As a test for now use the last destination
	d := dd[len(dd)-1]

	// Shutting it down deferred
	defer clientAPI.DestinationService().ShutDown(ctx, d)

	// Booting it now
	clientAPI.DestinationService().Boot(ctx, d)

	// Running the test
	cfg.Destination = d
	return clientAPI.ActionRunTest().Run(ctx, cfg)
}
*/
