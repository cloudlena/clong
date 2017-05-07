# Clong

A simple game where players have to throw balls at targets from their smart phones.

1. Open `/screen` on any big screen. This is where the game runs. The game should begin to spawn targets.
1. Open `/` on any touch device and swipe forward to launch balls at the targets.

## Run locally

1. Run `go build ./cmd/clong`
1. Execute the created binary and visit <http://localhost:8080>

## Build with Docker

1. Run `docker run --rm -v "${PWD}:/go/src/github.com/mastertinner/clong" -w /go/src/github.com/mastertinner/clong golang go build ./cmd/clong`

    To cross-compile for windows, use the `-e "GOOS=windows" -e "GOARCH=amd64"` flags (depending on your system, you might have to adjust `GOARCH`)

    To cross-compile for macOS, use the `-e "GOOS=darwin" -e "GOARCH=amd64"` flags (depending on your system, you might have to adjust `GOARCH`)

## Run on Cloud Foundry

1. Change `host` in `manifest.yml` to something that isn't taken yet
1. Run `GOOS=linux GOARCH=amd64 go build ./cmd/clong`
1. Run `cf push`
