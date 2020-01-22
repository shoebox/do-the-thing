package main

import (
	"dothething/internal/util"
	"dothething/internal/xcode"
	"fmt"
)

func main() {
	fmt.Println("main")

	/*
		res, err := service.List()
		fmt.Println(res, err)
		for _, r := range res {
			fmt.Println(r)
		}
	*/

	exec := util.OsExec{}
	fileUtilService := util.IoUtilFileService{}

	listService := xcode.NewXCodeListService(exec, fileUtilService)
	list, err := listService.List()
	fmt.Println(list, err)
	for _, v := range list {
		fmt.Println(v)
	}
	fmt.Println(listService.List())

	selectService := xcode.NewSelectService(listService, exec)

	res, err := selectService.SelectVersion("10.2.1")
	fmt.Println(res, err)
}
