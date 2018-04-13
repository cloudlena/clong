all:
	go build ./cmd/clong

test:
	go test ./... -race -cover

build-docker:
	docker build . -f build/docker/Dockerfile -t clong

deploy-cf:
	GOOS=linux GOARCH=amd64 go build ./cmd/clong
	cf push -f deployments/cf/manifest.yml
