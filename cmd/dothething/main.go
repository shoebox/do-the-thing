package main

import (
	"context"
	"dothething/internal/action/unittest"
	"dothething/internal/destination"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	path := "/Users/johann.martinache/Desktop/tmp/Swiftstraints/Swiftstraints.xcodeproj"

	executor := util.NewExecutor()
	xcodeService := xcode.NewService(executor, path)

	// pj := xcode.NewProjectService(xcodeService)

	dest := destination.NewDestinationService(xcodeService, executor)

	//
	dd, err := dest.List(ctx, "SwiftstraintsTests")
	fmt.Println(err)
	fmt.Println(dd)

	dest.Boot(ctx, dd[0])

	//
	a := unittest.NewActionRun(xcodeService, executor)
	err = a.Run(ctx, dd[0].Id)
	fmt.Println(err)

	defer dest.ShutDown(ctx, dd[0])

	/*
		// # find the id that points to the location of the encoded file in the .xcresult bundle
		// id=$(xcrun xcresulttool get --format json --path Tests.xcresult | jq '.actions._values[]' | jq -r '.actionResult.logRef.id._value')
		// # export the log found at the the id in the .xcresult bundle
		// xcrun xcresulttool export --path Tests.xcresult --id $id --output-path TestsStdErrorStdout.log --type file
	*/

	// p, err := pj.Parse()
	// fmt.Println(p, err)

	// dd, err := pj.ListDestinations("test")
	// fmt.Println(dd, err)

	/*
		k, err := keychain.NewKeyChain(executor)
		fmt.Println(k, err)

		err = k.Create(ctx, "password")
		if err != nil {
			fmt.Println(err)
		}

		err = k.ImportCertificate(ctx, "assets/Certificate.p12", "p4ssword", "123")
		if err != nil {
			fmt.Println(err)
		}
	*/

	/*
		// defer k.Delete()

		// f, err := k.Create()
		// fmt.Println(f.Name(), err)
	*/

	/*
		fus := util.IoUtilFileService{}

		listService := xcode.NewXCodeListService(executor, fus)
		list, err := listService.List(ctx)
		fmt.Println(list, err)
		for _, v := range list {
			fmt.Println(v)
		}

		selectService := xcode.NewSelectService(listService, executor)

		res, err := selectService.Find(ctx, "11.3.1")
		fmt.Println(res, err)
	*/
}
