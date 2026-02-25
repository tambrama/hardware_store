# ===== build stage =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY hardware_store/go.mod ./go.mod
COPY hardware_store/go.sum ./go.sum
RUN go mod download

COPY hardware_store/ ./ 
RUN go build -o app ./cmd/app

FROM alpine:3.19

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/internal/config /app/internal/config

EXPOSE 8081

CMD ["./app"]
