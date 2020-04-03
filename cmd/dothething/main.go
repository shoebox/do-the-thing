package main

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	exec := util.NewCommandRunner()

	path := "/Users/johann.martinache/Desktop/tmp/toto/test/test.xcodeproj"

	xcodeService := xcode.NewService(exec, path)

	pj := xcode.NewProjectService(xcodeService)
	fmt.Println(pj.Parse(ctx))

	dest := destination.NewDestinationService(xcodeService, exec)
	d, err := dest.List(ctx, "test")
	fmt.Println(d, err)

	err = dest.Boot(ctx, d[0])
	fmt.Println(err)

	err = dest.ShutDown(ctx, d[0])
	fmt.Println("shutdown ", err)

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
		k := keychain.NewKeyChain(exec)
		fmt.Println(k)

		err := k.Create("password")
		fmt.Println(err)

		f, err := os.Open("assets/Certificate.p12")
		fmt.Println(f, err)

		err = k.ImportCertificate("assets/Certificate.p12", "p4ssword", "123")
		fmt.Println(err)

		// defer k.Delete()

		// f, err := k.Create()
		// fmt.Println(f.Name(), err)

			listService := xcode.NewXCodeListService(exec, fileUtilService)
			list, err := listService.List()
			fmt.Println(list, err)
			for _, v := range list {
				fmt.Println(v)
			}
			fmt.Println(listService.List())

			selectService := xcode.NewSelectService(listService, exec)

			res, err := selectService.Find("10.2.1")
			fmt.Println(res, err)
	*/
}
