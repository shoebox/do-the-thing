package main

import (
	"dothething/internal/xcode"
	"fmt"
)

func main() {
	fmt.Println("main")
	service := xcode.New()
	res, err := service.List()
	fmt.Println(res, err)
	for _, r := range res {
		fmt.Println(r)
	}
}
