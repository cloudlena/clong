all:
	go build ./cmd/clong

build-docker:
	docker run --rm -v "${PWD}:/go/src/github.com/mastertinner/clong" -w /go/src/github.com/mastertinner/clong golang go build ./cmd/clong

deploy-cf:
	GOOS=linux GOARCH=amd64 go build ./cmd/clong
	cf push
