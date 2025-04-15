package api

import "os/exec"

func execCMD(AppConfig struct, cmd string) {
	command := []string{
		"-r",
		AppConfig.resticRepoPath,
	}

	exec.Command("restic")
}
