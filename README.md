# Clong

[![Go Report Card](https://goreportcard.com/badge/github.com/cloudlena/clong)](https://goreportcard.com/report/github.com/cloudlena/clong)
[![Build Status](https://github.com/cloudlena/clong/actions/workflows/main.yml/badge.svg)](https://github.com/cloudlena/clong/actions)

A simple game where players have to throw balls at targets from their smartphones.

1.  Open `/screen` on any big screen. This is where the game runs. The game should begin to spawn targets.
1.  Open `/` on any touch device and swipe forward to launch balls at the targets. Many players can play at the same time.
1.  Open `/scoreboard` to get a list of high scores (which updates live).

## Build and Run Locally

1.  Run `make`
1.  Execute the created binary and visit <http://localhost:8080>

## Run Tests

1.  Run `make test`

## Build Container Image

The image is also available on [Docker Hub](https://hub.docker.com/r/cloudlena/clong/).

1.  Run `make build-docker`

## Run on Cloud Foundry

1.  Create an SQL database service
1.  Modify `deployments/cf/*` to your liking
1.  Run `make deploy-cf`
