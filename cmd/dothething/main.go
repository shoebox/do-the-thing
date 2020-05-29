package main

import (
	"context"
	"dothething/internal/action/unittest"
	"dothething/internal/destination"
	"dothething/internal/keychain"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/project"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var xcb xcode.BuildService
var listService xcode.ListService
var pj project.ProjectService

func main() {
	// logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	//
	// path := "/Users/johann.martinache/Desktop/tmp/Swiftstraints/Swiftstraints.xcodeproj"
	path := "/Users/johann.martinache/Desktop/massive/bein/bein-apple/beIN.xcodeproj"

	f := util.IoUtilFileService{}
	e := util.NewExecutor()
	xcb = xcode.NewService(e, path)
	pj = project.NewProjectService(xcb, e)
	prj, err := pj.Parse(ctx)
	fmt.Printf("Project %#v %v\n", prj, err)
	fmt.Printf("Project %#v %v\n", prj.Schemes, err)

	// List service
	listService := xcode.NewXCodeListService(e, f)
	fmt.Println(listService)
	//listService.List(ctx)

	/*
		list, err := listService.List(ctx)
		fmt.Println("Xcode list :::", list, err)

		//
		if err := selectService(e, listService); err != nil {
			log.Error().AnErr("Error", err).Msg("Select service error")
		}

		//
		if err := unitTest(e, xcodeService); err != nil {
			log.Error().AnErr("Error", err).Msg("Unit test error")
		}

		if err := keychainTest(e); err != nil {
			log.Error().AnErr("Error", err).Msg("Keychain error")
		}
		/*
			// # find the id that points to the location of the encoded file in the .xcresult bundle
			// id=$(xcrun xcresulttool get --format json --path Tests.xcresult | jq '.actions._values[]' | jq -r '.actionResult.logRef.id._value')
			// # export the log found at the the id in the .xcresult bundle
			// xcrun xcresulttool export --path Tests.xcresult --id $id --output-path TestsStdErrorStdout.log --type file
	*/
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

func unitTest(e util.Executor, x xcode.BuildService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	dest := destination.NewDestinationService(x, e)
	dd, err := dest.List(ctx, "Swiftstraints iOS")
	if err != nil {
		return err
	}
	defer dest.ShutDown(ctx, dd[len(dd)-1])

	if err := dest.Boot(ctx, dd[len(dd)-1]); err != nil {
		return err
	}

	a := unittest.NewActionRun(x, e)
	return a.Run(ctx, dd[0].Id)
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
	defer k.Delete(ctx)

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
