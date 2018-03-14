FROM golang:latest

WORKDIR /go/src/github.com/mastertinner/clong
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go build ./cmd/clong

EXPOSE 8080

ENTRYPOINT ["./clong"]
