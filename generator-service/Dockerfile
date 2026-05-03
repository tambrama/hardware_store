# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Копируем go.mod и go.sum
COPY ./generator-service/go.mod ./generator-service/go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY generator-service/ ./

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o generator ./cmd/

# Stage 2: Runtime
FROM alpine:3.19

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем бинарник из builder'а
COPY --from=builder /build/generator .

# Копируем конфиг
COPY ./generator-service/config/ ./config/

# Переменная окружения для пути конфига
ENV CONFIG_PATH=/app/config/local.yaml

ENTRYPOINT ["./generator"]
