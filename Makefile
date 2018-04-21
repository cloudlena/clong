.PHONY: all lint test build-docker deploy-cf

all:
	go build ./cmd/...

lint:
	gometalinter --vendor ./...

test:
	go test -race -cover ./...

build-docker:
	docker build -f build/docker/Dockerfile -t clong .

deploy-cf:
	GOOS=linux go build -ldflags="-s -w" ./cmd/clong
	cf push -f deployments/cf/manifest.yml
