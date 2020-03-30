package main

import (
	"dothething/internal/destination"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	exec := util.NewCommandRunner()
	// fileUtilService := util.IoUtilFileService{}

	path := "/Users/johann.martinache/Desktop/tmp/toto/test/test.xcodeproj"

	xcodeService := xcode.NewService(exec, path)

	pj := xcode.NewProjectService(xcodeService)
	fmt.Println(pj.Parse())

	dest := destination.NewDestinationService(xcodeService, exec)
	d, err := dest.List("test")
	fmt.Println(d, err)

	err = dest.Boot(d[0])
	fmt.Println(err)

	err = dest.ShutDown(d[0])
	fmt.Println("shutdown ", err)

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
