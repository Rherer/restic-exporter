package main

import (
	"fmt"
	"os"
)

func main() {
	os.Setenv("RESTIC_PASSWORD", "1")
	os.Setenv("RESTIC_PASSWORD_FILE", "./pw.txt")
	os.Setenv("RESTIC_REPOSITORY", "./restic")
	resticConfig := config.getConfig()
	fmt.Println(resticConfig)
}
