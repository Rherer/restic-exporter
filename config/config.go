package config

import "os"

type config struct {
	resticPW       string
	resticPWFile   string
	resticRepoPath string
}

func getConfig() *config {
	resticConfig := config{
		resticPW:       os.Getenv("RESTIC_PASSWORD"),
		resticPWFile:   os.Getenv("RESTIC_PASSWORD_FILE"),
		resticRepoPath: os.Getenv("RESTIC_REPOSITORY"),
	}

	return &resticConfig
}
