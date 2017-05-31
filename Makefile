all:
	go build ./cmd/clong

deploy-cf:
	GOOS=linux GOARCH=amd64 go build ./cmd/clong
	cf push
