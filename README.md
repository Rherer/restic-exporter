# Restic Prometheus Exporter
![Go Version][go-vers]
![Go Ranking][go-rank]
![Build Status][build-stat]
![Downloads][github-downloads]
![Contributors][github-contribs]

> Exports statistics about a Restic Repository as prometheus metrics
> Metrics are refreshed on each scrape, checks will (optionally) be run in a configured interval

This project was inspired by: https://github.com/ngosang/restic-exporter (so props to @ngosang)
The metrics are not 1:1 replaceable, but i did my best to keep compatibility high.

## Installation

 There are multiple ways to install and run this program:
> - Binary
> - Systemd Unit
> - Container

You can either download a prebuilt binary from the releases tab (provided for Windows and Linux), or build your own using the following steps.
> For Linux a optional, automated setup script is provided, which:
> - builds (or downloads, if go is not installed) the binary
> - optionally installs a systemd-service

Alternatively you can just run the provided docker container

### Docker:

```sh
docker run -d \
-e RESTIC_REPOSITORY='/etc/restic-repo' \
-e RESTIC_PASSWORD_FILE='/etc/restic-pw' \
--mount type=bind,src=/foo/bar/,dst=/etc/restic-repo \
--mount type=bind,src=/baz/qux,dst=/etc/restic-pw \
-p 80:8080 \
--name restic-exporter ghcr.io/rherer/restic-exporter:latest
```

### Linux, Windows & Mac OSX:

```sh
git clone https://github.com/Rherer/restic-exporter .
go mod download
go install
```

### Linux only:

```sh
git clone https://github.com/Rherer/restic-exporter .
sh install.sh
```

## Usage example

Configuration options (Using Environment Variables):

```sh
HTTP_BASE_PATH  The path that the metrics will be served on (Default: /metrics)
HTTP_BASE_PORT  The port on which the server will listen (Default: 8080)
CHECK_INTERVAL  The interval between checks (Default: 30m)
NO_CHECK        Disable periodic checks completely (also disables the corresponding metric) (Default: false)
USE_REPO_PATH   Add the path to the repository as an additional tag (Default: false)
USE_SNAPSHOT_ID Add the individual snapshot id as an additional tag (Default: false)
USE_LATEST_N    Collect the latest n snapshots (grouped by host and path) the repo (Default: 1)
```

All other options will be passed through as is, so you can natively use any options that restic provides as environment variables.

The following stats are currently supported:
```sh
- Per Snapshot:
  - Start time
  - Duration
  - No. of files
  - No. of new files
  - No. of modified files
  - Size in bytes
- Globally:
  - No. of locks
  - Status of last check (Configurable interval)
  - Total No. of snapshots
```

## Setup on Linux

> I recommend running the exporter as a systemd-unit or a docker container.
> You can find a minimal configuration for either method below.

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

You can then access the metrics in your browser at ex.: http://localhost:8080/metrics
```
# HELP restic_backup_files_changed Shows the amount of changed files in the snapshot
# TYPE restic_backup_files_changed gauge
restic_backup_files_changed{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="c"} 590
restic_backup_files_changed{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 883
# HELP restic_backup_files_new Shows the amount of new files in the snapshot
# TYPE restic_backup_files_new gauge
restic_backup_files_new{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="c"} 2372
restic_backup_files_new{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 38588
# HELP restic_backup_files_total Shows the total amount of files in the snapshot
# TYPE restic_backup_files_total gauge
restic_backup_files_total{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="c"} 33963
restic_backup_files_total{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 62937
# HELP restic_backup_runtime Shows the time the snapshot took
# TYPE restic_backup_runtime gauge
restic_backup_runtime{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="c"} 171
restic_backup_runtime{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 115
# HELP restic_backup_size_total Shows the amount of bytes in the snapshot
# TYPE restic_backup_size_total gauge
restic_backup_size_total{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="c"} 1.02404161055e+11
restic_backup_size_total{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 4.0023511166e+10
# HELP restic_backup_timestamp Shows the start time of the snapshot
# TYPE restic_backup_timestamp gauge
restic_backup_timestamp{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="c"} 1.747670774e+09
restic_backup_timestamp{client_hostname="fedora.fritz.box",client_username="eherer",client_version="restic 0.17.3",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 1.744825025e+09
# HELP restic_check_success Shows whether a check was successful
# TYPE restic_check_success gauge
restic_check_success 0
# HELP restic_locks_total Shows the amount of locks on the repository
# TYPE restic_locks_total gauge
restic_locks_total 0
# HELP restic_scrape_duration_seconds Shows the duration of the scrape
# TYPE restic_scrape_duration_seconds gauge
restic_scrape_duration_seconds 1.605832159
# HELP restic_snapshots_total Shows the total amount of snapshots in the repository
# TYPE restic_snapshots_total gauge
restic_snapshots_total{client_hostname="fedora.fritz.box",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="a,b"} 1
restic_snapshots_total{client_hostname="fedora.fritz.box",snapshot_paths="/mnt/data/distrobox/restic-exporter/home",snapshot_tags="c"} 1
restic_snapshots_total{client_hostname="fedora.fritz.box",snapshot_paths="/run/host/var/mnt/data/distrobox/restic-exporter/home",snapshot_tags=""} 2
```

You can now scrape the metrics endpoint using prometheus.

## Setup on Windows

> As running as a service under windows is not natively supported, either use
> - docker
> - or [nssm](https://nssm.cc/)

## Local Development setup

If building yourself, you can just use the same procedure as described above.:

```sh
go mod download
go build
```

## Meta

Distributed under the GPL 3.0 license. See ``LICENSE`` for more information.

## Contributing

1. Fork it (<https://github.com/Rherer/restic-exporter/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

<!-- Markdown link & img dfn's -->
[go-vers]: https://img.shields.io/github/go-mod/go-version/Rherer/restic-exporter
[go-rank]: https://goreportcard.com/badge/github.com/Rherer/restic-exporter
[build-stat]: https://img.shields.io/github/actions/workflow/status/Rherer/restic-exporter/build-release-binaries.yaml
[github-downloads]: https://img.shields.io/github/downloads/Rherer/restic-exporter/total
[github-contribs]: https://img.shields.io/github/contributors/Rherer/restic-exporter
