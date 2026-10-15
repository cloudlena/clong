FROM docker.io/library/golang:1 AS builder
WORKDIR /usr/src/app
COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/clong

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /usr/src/app/bin/clong /clong
EXPOSE 8080
ENTRYPOINT [ "/clong" ]
