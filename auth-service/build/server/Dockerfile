FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY ["../go.mod", "../go.sum", "./"]
RUN go mod download

COPY . .

RUN go install -tags sqlite3 github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o app ./cmd/auth

FROM alpine:3.19

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

COPY --from=builder /app/app .
COPY --from=builder /app/config /app/config

COPY --from=builder /app/migrations /app/migrations

COPY --from=builder /go/bin/goose /bin/goose

RUN mkdir -p /app/storage && chmod 777 /app/storage

EXPOSE 9090

CMD ["./app"]
