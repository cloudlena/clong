# Clong

[![Go Report Card](https://goreportcard.com/badge/github.com/mastertinner/clong?style=flat-square)](https://goreportcard.com/report/github.com/mastertinner/clong)
[![Build Status](https://img.shields.io/travis/mastertinner/clong.svg?style=flat-square)](https://travis-ci.org/mastertinner/clong)
[![Docker Build](https://img.shields.io/docker/build/mastertinner/clong.svg?style=flat-square)](https://hub.docker.com/r/mastertinner/clong)

A simple game where players have to throw balls at targets from their smart phones.

1.  Open `/screen` on any big screen. This is where the game runs. The game should begin to spawn targets.
1.  Open `/` on any touch device and swipe forward to launch balls at the targets. Many players can play at the same time.
1.  Open `/scoreboard` to get a list of high scores (which updates live).

## Build and Run Locally

1.  Run `make`
1.  Execute the created binary and visit <http://localhost:8080>

## Run Tests

1.  Run `make test`

## Build Docker Image

The image is available on [Docker Hub](https://hub.docker.com/r/mastertinner/clong/)

1.  Run `make build-docker`

## Run on Cloud Foundry

1.  Create an SQL database service
1.  Modify `deployments/cf/*` to your liking
1.  Run `make deploy-cf`
