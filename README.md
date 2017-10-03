# Clong

[![Go Report Card](https://goreportcard.com/badge/github.com/mastertinner/clong)](https://goreportcard.com/report/github.com/mastertinner/clong)

A simple game where players have to throw balls at targets from their smart phones.

1. Open `/screen` on any big screen. This is where the game runs. The game should begin to spawn targets
1. Open `/` on any touch device and swipe forward to launch balls at the targets. Many players can play at the same time.
1. Open `/scoreboard` to get a list of high scores (which updates live)

## Run the app

1. Install [Dep](https://github.com/golang/dep)
1. Run `dep ensure`

### Run locally

1. Run `make`
1. Run a local database (e.g. with `docker run -d -p "3306:3306" -e "MYSQL_ALLOW_EMPTY_PASSWORD=yes" -e "MYSQL_DATABASE=clong" mysql`)
1. Execute the created binary and visit <http://localhost:8080>

### Run on Cloud Foundry

1. Change `host` in `manifest.yml` to something that isn't taken yet
1. Run `cf create-service mariadbent usage clong-db` to create a DB service if it doesn't exist yet
1. Run `make deploy-cf`
