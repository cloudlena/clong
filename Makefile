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

.PHONY: clean
clean:
	rm -rf bin
