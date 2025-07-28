.PHONY: build
build:
	go build -o bin/clong

.PHONY: run
run:
	go run main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: build-image
build-image:
	podman build -t clong .

.PHONY: deploy-cf
deploy-cf:
	GOOS=linux go build -ldflags="-s -w" -o bin/clong ./cmd/clong
	cf push -f deployments/cf/manifest.yml

.PHONY: clean
clean:
	rm -rf bin
