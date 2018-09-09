FROM golang:1 as builder
RUN groupadd -r clong && useradd --no-log-init -r -g clong clong
WORKDIR /usr/src/clong
COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a -installsuffix cgo -o bin/clong ./cmd/clong

FROM scratch
WORKDIR /usr/clong
COPY --from=builder /usr/src/clong/bin/clong /usr/src/clong/web ./
COPY --from=builder /etc/passwd /etc/passwd
USER clong
EXPOSE 8080
ENTRYPOINT ["./clong"]
