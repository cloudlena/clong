.PHONY: all lint test build-docker deploy-cf clean

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

all:
	go build -o bin/clong ./cmd/clong

lint:
	golangci-lint run

test:
	go test -race -cover ./...

build-docker:
	docker build -t clong .

deploy-cf:
	GOOS=linux go build -ldflags="-s -w" -o bin/clong ./cmd/clong
	cf push -f deployments/cf/manifest.yml

clean:
	rm -rf bin
