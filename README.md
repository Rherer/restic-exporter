# Restic Prometheus Exporter
![Go Version][go-vers]
![Build Status][build-stat]
![Downloads][github-downloads]
![Contributors][github-contribs]

> Exports statistics about a Restic Repository as prometheus metrics
> Metrics are refreshed on each scrape, checks will (optionally) be run in a configured interval
> This project was inspired by: https://github.com/ngosang/restic-exporter (so props to @ngosang)
> The metrics are not 1:1 replaceable, but i did my best to keep compatibility high

> The following metrics are exported:
>  - Per Snapshot:
>    - Start time
>    - Duration
>    - No. of files
>    - No. of new files
>    - No. of modified files
>    - Size in bytes
>  - Globally:
>    - No. of locks
>    - Status of last check (Configurable interval)
>    - Total No. of snapshots

## Installation

You can download a prebuilt binary from the releases tab, or build your own using the following steps.
For Linux an automated setup script is provided, which:
- builds (or downloads, if go is not installed) the binary
- optionally installs a systemd-service

Linux, Windows & Mac OSX:

```sh
git clone https://github.com/Rherer/restic-exporter .
go mod download
go install
```

Linux only:

```sh
git clone https://github.com/Rherer/restic-exporter .
sh install.sh
```

## Usage example

Configuration options (Using Environment Variables):

```sh
	HTTP_BASE_PATH The path that the metrics will be served on (Default: /metrics)
	HTTP_BASE_PORT The port on which the server will listen (Default: 8080)
	CHECK_INTERVAL The interval between checks (Default: 30m)
	NO_CHECK       Disable periodic checks completely (also disables the corresponding metric) (Default: false)
	USE_REPO_PATH  Add the path to the repository as an additional tag (Default: false)
```

### Setup on Linux

First you have to set the configuration for restic with your environment.
The options will be passed through as is, so you can use any and all options that restic provides as environment variables.
You can set these variables before calling the binary
or add them to the systemd-unit <-- This is the recommended way!

Example for a minimal systemd unit-file:
```sh
[Unit]
Description=Prometheus Exporter for Restic
After=network.target

[Service]
User=restic
ENVIRONMENT=RESTIC_REPOSITORY=/mnt/backup/restic
ENVIRONMENT=RESTIC_PASSWORD_FILE=/mnt/backup/restic_pw
Restart=on_fail
ExecStart=/usr/bin/restic-exporter

[Install]
WantedBy=default.target
```

You can then access the metrics in your browser like:

http://localhost:8080/metrics
```
# HELP restic_backup_files_changed Shows the amount of changed files in the snapshot
# TYPE restic_backup_files_changed gauge
restic_backup_files_changed{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="4544e9a1740a21a800c69603b959c2189545ffae254db2557ef6b26b6835c8cb",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 883
restic_backup_files_changed{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="d0b273a4bd2fab20558b9a0e1fa8aece2645d8e33bdb9fb7fd5fe7df3de123cd",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 0
# HELP restic_backup_files_new Shows the amount of new files in the snapshot
# TYPE restic_backup_files_new gauge
restic_backup_files_new{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="4544e9a1740a21a800c69603b959c2189545ffae254db2557ef6b26b6835c8cb",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 38588
restic_backup_files_new{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="d0b273a4bd2fab20558b9a0e1fa8aece2645d8e33bdb9fb7fd5fe7df3de123cd",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 24359
# HELP restic_backup_files_total Shows the total amount of files in the snapshot
# TYPE restic_backup_files_total gauge
restic_backup_files_total{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="4544e9a1740a21a800c69603b959c2189545ffae254db2557ef6b26b6835c8cb",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 62937
restic_backup_files_total{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="d0b273a4bd2fab20558b9a0e1fa8aece2645d8e33bdb9fb7fd5fe7df3de123cd",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 24359
# HELP restic_backup_runtime Shows the time the snapshot took
# TYPE restic_backup_runtime gauge
restic_backup_runtime{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="4544e9a1740a21a800c69603b959c2189545ffae254db2557ef6b26b6835c8cb",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 115
restic_backup_runtime{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="d0b273a4bd2fab20558b9a0e1fa8aece2645d8e33bdb9fb7fd5fe7df3de123cd",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 105
# HELP restic_backup_size_total Shows the amount of bytes in the snapshot
# TYPE restic_backup_size_total gauge
restic_backup_size_total{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="4544e9a1740a21a800c69603b959c2189545ffae254db2557ef6b26b6835c8cb",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 4.0023511166e+10
restic_backup_size_total{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="d0b273a4bd2fab20558b9a0e1fa8aece2645d8e33bdb9fb7fd5fe7df3de123cd",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 1.9628940161e+10
# HELP restic_backup_timestamp Shows the start time of the snapshot
# TYPE restic_backup_timestamp gauge
restic_backup_timestamp{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="4544e9a1740a21a800c69603b959c2189545ffae254db2557ef6b26b6835c8cb",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 1.744825025e+09
restic_backup_timestamp{client_hostname="fedora",client_username="Rherer",client_version="restic 0.17.3",snapshot_id="d0b273a4bd2fab20558b9a0e1fa8aece2645d8e33bdb9fb7fd5fe7df3de123cd",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 1.744735999e+09
# HELP restic_check_success Shows whether a check was sucessfull
# TYPE restic_check_success gauge
restic_check_success 0
# HELP restic_locks_total Shows the amount of locks on the repository
# TYPE restic_locks_total gauge
restic_locks_total 0
# HELP restic_scrape_duration_seconds Shows the duration of the scrape
# TYPE restic_scrape_duration_seconds gauge
restic_scrape_duration_seconds 1.043720226
# HELP restic_snapshots_total Shows the total amount of snapshots in the repository
# TYPE restic_snapshots_total gauge
restic_snapshots_total 2
```

You can then scrape the metrics using prometheus, and display them using grafana.

### Setup on Windows

Add your needed environment variables in your settings.

Then create a service from the binary
> You can use nssm for running a binary as a service under windows
> See the docs under: https://nssm.cc/

You can then scrape the metrics using prometheus, and display them using grafana.

## Development setup

```sh
go mod download
go build
```

## Meta

Distributed under the GPL 3.0 license. See ``LICENSE`` for more information.

[https://github.com/yourname/github-link](https://github.com/dbader/)

## Contributing

1. Fork it (<https://github.com/Rherer/restic-exporter/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

<!-- Markdown link & img dfn's -->
[go-vers]: https://img.shields.io/github/go-mod/go-version/Rherer/restic-exporter
[build-stat]: https://img.shields.io/github/actions/workflow/status/Rherer/restic-exporter/build-release-binaries.yaml
[github-downloads]: https://img.shields.io/github/downloads/Rherer/restic-exporter/total
[github-contribs]: https://img.shields.io/github/contributors/Rherer/restic-exporter
