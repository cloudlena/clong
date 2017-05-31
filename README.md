# Clong

A simple game where players have to throw balls at targets from their smart phones.

1. Open `/screen` on any big screen. This is where the game runs. The game should begin to spawn targets.
1. Open `/` on any touch device and swipe forward to launch balls at the targets.

## Run the app

1. Install (Dep)[https://github.com/golang/dep]
1. Run `dep ensure`

### Run locally

1. Run `make`
1. Execute the created binary and visit <http://localhost:8080>

### Run on Cloud Foundry

1. Change `host` in `manifest.yml` to something that isn't taken yet
1. Run `make deploy-cf`
