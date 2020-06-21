package main

import (
	"context"
	"dothething/internal/client"
	"dothething/internal/config"
	"dothething/internal/keychain"
	"dothething/internal/signature"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/project"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfg config.Config
var xcb xcode.BuildService
var fileUtil util.FileService
var serviceList xcode.ListService

//var serviceProject project.ProjectService
var serviceProvisioning signature.ProvisioningService
var executor util.Executor
var targetProject project.Project

var api client.API

func main() {
	// logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cfg = config.Config{
		Path:          "/Users/johann.martinache/Desktop/tmp/BookStore-iOS/BookStore.xcodeproj",
		Scheme:        "BookStore",
		Configuration: "Release",
		Target:        "BookStore",
		CodeSignOption: config.SignConfig{
			Path: "",
		},
	}

	var err error
	api, err = client.NewAPIClient(cfg)
	if err != nil {
		panic(err)
	}

	err = resolveSignature(ctx)
	fmt.Println("err", err)
	build()
	archive()
	unitTest()
}

func listXCodeInstance(ctx context.Context) error {
	installs, err := api.XCodeListService().List(ctx)
	if err != nil {
		return err
	}

	for _, i := range installs {
		fmt.Println(i.Path, i.Version)
	}

	return nil
}

func resolveSignature(ctx context.Context) error {
	pj, err := api.XCodeProjectService().Parse(ctx)
	if err != nil {
		return err
	}

	if err := pj.ValidateConfiguration(cfg); err != nil {
		log.Panic().AnErr("Error", err)

		return err
	}

	res, err := api.SignatureResolver().Resolve(ctx, cfg, pj)
	if err != nil {
		return err
	}
	fmt.Println(res, err)
	return nil
}

func selectService(e util.Executor, l xcode.ListService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	selectService := xcode.NewSelectService(l, e)
	res, err := selectService.Find(ctx, "11.3.1")
	if err != nil {
		return err
	}

	fmt.Println("XCode instance : ", res, err)
	return nil
}

func build() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	return api.
		ActionBuild().
		Run(ctx, cfg)
}

func archive() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	return api.
		ActionArchive().
		Run(ctx, cfg)
}

func unitTest() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Retrieving the destination for the scheme
	dd, err := api.
		DestinationService().
		List(ctx, cfg.Scheme)

	if err != nil {
		return err
	}

	// As a test for now use the last destination
	d := dd[len(dd)-1]

	// Shutting it down deferred
	defer api.DestinationService().ShutDown(ctx, d)

	// Booting it now
	api.DestinationService().Boot(ctx, d)

	// Running the test
	return api.ActionRunTest().Run(ctx, d, cfg)
}

func keychainTest(e util.Executor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Keychain service
	k, err := keychain.NewKeyChain(e)
	if err != nil {
		return err
	}

	// Delete it at the end
	// defer k.Delete(ctx)

	// Create the new keychain
	err = k.Create(ctx, "password")
	if err != nil {
		return err
	}

	// Import the certificate
	err = k.ImportCertificate(ctx, "assets/Certificate.p12", "p4ssword", "123")
	if err != nil {
		return err
	}

	return nil
}
