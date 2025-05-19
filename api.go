package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Snapshot struct {
	Time           time.Time `json:"time"`
	Tree           string    `json:"tree"`
	Parent         string    `json:"parent"`
	Paths          []string  `json:"paths"`
	Tags           []string  `json:"tags"`
	Hostname       string    `json:"hostname"`
	Username       string    `json:"username"`
	UID            int       `json:"uid"`
	Gid            int       `json:"gid"`
	ProgramVersion string    `json:"program_version"`
	Summary        struct {
		BackupStart         time.Time `json:"backup_start"`
		BackupEnd           time.Time `json:"backup_end"`
		FilesNew            int       `json:"files_new"`
		FilesChanged        int       `json:"files_changed"`
		FilesUnmodified     int       `json:"files_unmodified"`
		DirsNew             int       `json:"dirs_new"`
		DirsChanged         int       `json:"dirs_changed"`
		DirsUnmodified      int       `json:"dirs_unmodified"`
		DataBlobs           int       `json:"data_blobs"`
		TreeBlobs           int       `json:"tree_blobs"`
		DataAdded           int64     `json:"data_added"`
		DataAddedPacked     int64     `json:"data_added_packed"`
		TotalFilesProcessed int       `json:"total_files_processed"`
		TotalBytesProcessed int64     `json:"total_bytes_processed"`
	} `json:"summary"`
	ID      string `json:"id"`
	ShortID string `json:"short_id"`
}

func checkIfRepoExists() {
	_, err := execCMD([]string{"cat", "config"})
	if err != nil {
		panic(err)
	}
}

// Get Snapshots in the repo
func getSnapshots() ([]Snapshot, error) {
	rawJson, err := execCMD([]string{"snapshots", "--latest", strconv.Itoa(Config.USE_LATEST_N)})
	if err != nil {
		return nil, err
	}

	var snapshots []Snapshot
	err = json.Unmarshal(rawJson, &snapshots)
	if err != nil {
		return snapshots, err
	}

	return snapshots, err
}

// More stats can probably be gotten from this
func getSnapshotCount() (int, error) {
	rawJson, err := execCMD([]string{"snapshots"})
	if err != nil {
		return 0, err
	}

	var snapshots []Snapshot
	err = json.Unmarshal(rawJson, &snapshots)
	if err != nil {
		return len(snapshots), err
	}

	return len(snapshots), err
}

// Run check, and return the exit code
func runCheck() (int, error) {
	_, err := execCMD([]string{"check"})
	if err != nil {
		//fmt.Println(result)
		return err.(*exec.ExitError).ExitCode(), err
	}

	return 0, err
}

func getLocks() (int, error) {
	result, err := execCMD([]string{"list", "locks"})
	if err != nil {
		return 0, err
	}

	return bytes.Count(result, []byte{'\n'}), err
}

// Execute passed restic command, and return the output
func execCMD(cmdString []string) ([]byte, error) {
	args := []string{
		"--json",
		"--no-lock",
	}

	args = append(args, cmdString...)

	cmd := exec.Command("restic", args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return output, err
}
