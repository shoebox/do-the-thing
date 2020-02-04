package main

import (
	"dothething/internal/util"
	"dothething/internal/xcode"
	"fmt"
)

func main() {
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

	res, err := selectService.Find("10.2.1")
	fmt.Println(res, err)

	res2 := xcode.NewProjectService(exec)
	err = res2.Parse()
	fmt.Println("err", err, res2)
}
