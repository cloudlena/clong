FROM golang:1 AS builder
RUN groupadd -r app && useradd --no-log-init -r -g app app
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a -installsuffix cgo -o bin/clong ./cmd/clong

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin/clong /app/web ./
COPY --from=builder /etc/passwd /etc/passwd
USER app
EXPOSE 8080
ENTRYPOINT ["./clong"]
