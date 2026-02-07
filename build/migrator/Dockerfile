FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY migrations ./migrations

# goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# ENTRYPOINT ["goose"]