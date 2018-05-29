.PHONY: all lint test build-docker deploy-cf

all:
	go build ./cmd/clong

lint:
	golangci-lint run --tests --enable-all

test:
	go test -race -cover ./...

build-docker:
	docker build -f build/docker/Dockerfile -t clong .

deploy-cf:
	GOOS=linux go build -ldflags="-s -w" ./cmd/clong
	cf push -f deployments/cf/manifest.yml
