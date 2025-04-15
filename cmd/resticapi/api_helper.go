package main

import "os/exec"

func execCMD(cmd string) {
	command := []string{
		"-r",
		AppConfig.resticRepoPath,
	}

	exec.Command("restic")
}
