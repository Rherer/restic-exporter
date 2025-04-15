package main

import (
	"fmt"
	"os"
)

type Config struct {
	resticPW       string `yaml: "restic_pw", envconfig: "RESTIC_PASSWORD"`
	resticPWFile   string `yaml: "restic_pw_file", envconfig: "RESTIC_PASSWORD_FILE"`
	resticRepoPath string `yaml: "restic_repo", envconfig: "RESTIC_REPOSITORY"`
}

var AppConfig Config

func main() {
	//AppConfig.resticPW = "1"
	AppConfig.resticPWFile = "/mnt/data/distrobox/restic-exporter/home/src/pw.txt"
	AppConfig.resticRepoPath = "../restic"

	validateConfig()

	resticapi.execCMD()
	fmt.Println(AppConfig)
}

func validateConfig() {
	if AppConfig.resticPWFile != "" {
		if AppConfig.resticPW != "" {
			panic("Both a password and a password file are configured. Remove either!")
		} else {
			buf, err := os.ReadFile(AppConfig.resticPWFile)
			if err != nil {
				panic(err)
			}
			AppConfig.resticPW = string(buf)
		}
	}

}
