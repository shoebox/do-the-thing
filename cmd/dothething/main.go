package main

import (
	"dothething/internal/keychain"
	"dothething/internal/util"
	"fmt"
	"os"
)

func main() {
	exec := util.OsExec{}
	// fileUtilService := util.IoUtilFileService{}

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

	/*
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
