# Clong

[![Go Report Card](https://goreportcard.com/badge/github.com/cloudlena/clong)](https://goreportcard.com/report/github.com/cloudlena/clong)
[![Build Status](https://github.com/cloudlena/clong/actions/workflows/main.yml/badge.svg)](https://github.com/cloudlena/clong/actions)

A multiplayer game where players have to throw balls at targets from their smartphones.

1. Open `/screen` on any big screen and log in as `admin` with your admin password. This is where the game runs. The game should begin to spawn targets.
1. Open `/` on any touch device and swipe forward to launch balls at the targets. Many players can play at the same time.
1. Open `/scoreboard` to get a list of high scores (which updates live).

## Configuration

Clong is configured through the following environment variables:

- `ADMIN_PASSWORD` (required): Password of the `admin` user, which is needed to open `/screen` and to reset scores
- `DATABASE_URL`: URL of the PostgreSQL database (defaults to `postgresql://postgres:clong@?sslmode=disable`)
- `PORT`: Port to listen on (defaults to `8080`)

## Resetting Scores

Run the following command to reset the scoreboard (replacing `PASSWORD` with your actual admin password):

```shell
curl -X DELETE localhost:8080/api/scores -u 'admin:PASSWORD'
```

## Run Locally

1. Run `docker compose up`
1. Visit <http://localhost:8080> (the admin password is `clong` unless you set `ADMIN_PASSWORD`)

## Build and Run Binary

1. Run `make`
1. Start a PostgreSQL database and set `DATABASE_URL` accordingly if it doesn't match the default
1. Run `ADMIN_PASSWORD=PASSWORD bin/clong` and visit <http://localhost:8080>

## Run Tests

1. Run `make test`

## Build Container Image

The image is also available on [Docker Hub](https://hub.docker.com/r/cloudlena/clong/).

1. Run `make build-image`

## Run on Kubernetes

1. Create a namespace and target it.
1. Define a PASSWORD for the `admin` user.
1. Define a DB_USERNAME and a DB_PASSWORD for clong to access the DB with. Stick to letters and digits, since they end up in the database URL.
1. Create a secret called `clong-credentials` as follows:

   ```shell
   kubectl create secret generic clong-credentials --from-literal=adminPassword="${PASSWORD}" --from-literal=dbUsername="${DB_USERNAME}" --from-literal=dbPassword="${DB_PASSWORD}"
   ```

1. Replace the host `clong.local` in `deployments/k8s/ing-clong.yml` with your own host.
1. Apply the deployment as follows:

   ```shell
   kubectl apply -f deployments/k8s
   ```

The database doesn't use a persistent volume, so scores are lost when its pod restarts.

## Run on Fly

1. Run `fly launch --config deployments/fly/fly.toml`
1. Create a PostgreSQL database and set the `DATABASE_URL` and `ADMIN_PASSWORD` secrets with `fly secrets set`
