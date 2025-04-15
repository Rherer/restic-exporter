package main

import (
	"fmt"
)

func main() {

	config.validateConfig()

	api.execCMD()
	fmt.Println(AppConfig)
}
