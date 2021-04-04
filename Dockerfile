FROM docker.io/library/golang:1 AS builder
WORKDIR /usr/src/app
COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a -installsuffix cgo -o bin/clong

FROM docker.io/library/alpine:latest
WORKDIR /usr/src/app
RUN addgroup -S clong && adduser -S clong -G clong
RUN apk add --no-cache dumb-init
COPY --from=builder --chown=clong:clong /usr/src/app/bin/clong /usr/src/app/web ./
USER clong
EXPOSE 8080
ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]
CMD [ "/usr/src/app/clong" ]
