# Clong

A simple game where players have to throw balls at targets from their smart phones.

1. Open `/screen` on any big screen. This is where the game runs. The game should begin to spawn targets.
1. Open `/` on any touch device and swipe forward to launch balls at the targets.

## Run locally

1. Run `make`
1. Execute the created binary and visit <http://localhost:8080>

## Build with Docker

1. Run `make build-docker`

    To cross-compile for windows, add the `-e "GOOS=windows" -e "GOARCH=amd64"` flags to the `Makefile` (depending on your system, you might have to adjust `GOARCH`)

    To cross-compile for macOS, add the `-e "GOOS=darwin" -e "GOARCH=amd64"` flags to the `Makefile` (depending on your system, you might have to adjust `GOARCH`)

## Run on Cloud Foundry

1. Change `host` in `manifest.yml` to something that isn't taken yet
1. Run `make deploy-cf`
