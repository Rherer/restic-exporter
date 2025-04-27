#!/bin/sh
# Build (or download) the binary, install it to /usr/bin and (optionally) set up a systemd-service

set -eu

GIT_LINK=https://github.com/Rherer/restic-exporter/releases/latest/restic-exporter_x64_linux

function start () {

    if ! command -v go 2>&1 >/dev/null; then
        echo "Go is not installed, trying to install pre-built binary... Continue?"
        select yn in "Yes" "No"; do
            case $yn in
                Yes ) download_bin; break;;
                No ) echo "Aborted."; exit;;
            esac
        done
    fi

    build
}

# Download Binary from releases if no golang is installed
function download_bin() {
    if [[ $(uname -m) != "x86_64" ]]; then
        echo "Unsupported architecture for prebuilt binary."
        echo "Please install golang and rerun this script."
        exit 1
    fi

    curl -L https://github.com/Rherer/restic-exporter/releases/latest/restic-exporter_x64_linux -o ./restic-exporter

    install_binary
}

# Build the binary for the current system on demand
function build() {
    go mod download
    go build -ldflags "-s" -o ./restic-exporter

    install_binary
}

# Install the binary to /usr/bin/
function install_binary(){
    mv ./restic-exporter /usr/bin/restic-exporter
    echo "Binary installed successfully."

    echo "Should a systemd unit be set up as well?"
    select yn in "Yes" "No"; do
        case $yn in
            Yes ) install_systemd_unit; break;;
            No ) echo "Aborted."; exit;;
        esac
    done
}

# Install a systemd-unit from the .service file
function install_systemd_unit() {
    if [ ! -f ./restic-exporter.service ]; then
        echo "Unit file not available, aborting..."
        exit 2
    fi

    echo "Please adjust the paths and Username in the .service file now, if not already done."
    echo "Continue?"
    select yn in "Yes" "No"; do
        case $yn in
            Yes ) break;;
            No ) echo "Aborted."; exit;;
        esac
    done

    cp ./restic-exporter.service /etc/systemd/system/restic-exporter.service 
    systemctl daemon-reload

    echo "Unit installed successfully!"
    echo "Enable automatic restart with: systemctl daemon enable restic-exporter.service"
    echo "Start with: systemctl daemon start restic-exporter.service"
}

start