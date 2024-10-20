FROM golang:1.22.0 AS builder

RUN go env -w GOMODCACHE=/root/.cache/go-build
RUN go env -w CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY / ./

RUN go build -o bin/app

FROM alpine:3.19.1

RUN apk update && apk upgrade

RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

RUN adduser -D appuser

USER appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /app/bin/ .
