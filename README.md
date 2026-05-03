## 🏪 Hardware Store (Go + Python) + Auth Service (gRPC SSO)

Монорепозиторий, содержащий:
- **hardware_store** — REST API магазина на Go (Gin, PostgreSQL).
- **store_python** — REST API магазина на Python (Flask, PostgreSQL).
- **photo-service** — отдельный сервис хранения фото (FastAPI) со шардированием на 4 базы Postgres.
- **auth-service** — отдельный сервис авторизации и управления сессиями (gRPC, JWT, SQLite).
- **generator-service** - отдельный сервис, который генерирует события об изменениях товаров и отправляет их в Kafka. 
- **protos** — общий модуль с gRPC контрактами (`.proto`) и сгенерированным кодом.

---

## 📂 Общая структура репозитория

```text
.
├── hardware_store/   # Go backend магазина (Gin)
├── store_python/     # Python Flask backend магазина
├── photo_service/    # Photo service (FastAPI) + sharding
├── auth-service/     # gRPC‑сервис авторизации (SSO)
├── generator-service/# Генератор событий
├── protos/           # Общие gRPC контракты и генерация кода
├── nginx/            # Конфигурация Nginx
├── docker/           # Init-скрипты для PostgreSQL
├── static/           # Статические файлы (index.html, image.png)
└── docker-compose.yml
```

Подробнее по каждому модулю:

- `store_python/README.md` — документация по Python-магазину.
- `protos/README.md` — описание gRPC контрактов, генерация кода для Go/Python.


---

## 🚀 Быстрый старт (Docker Compose)

В папке `src/`:

```bash
cd src
docker-compose up --build
```

После успешного старта:

- **API Go (hardware_store)**: `http://localhost/api/v1/` 
- **Swagger Go**: `http://localhost/swagger/` 
- **API Python (store_python)**: `http://localhost/api/v2/` 
- **Swagger Python**: `http://localhost/api/v2/apidocs/` 
- **SSE echo client**: `http://localhost/sse.html`
- **WebSocket echo client**: `http://localhost/ws.html`
- **pgAdmin**: `http://localhost/admin`
- **Nginx status**: `http://localhost/status` 
- **Auth gRPC сервис**: порт 9090 (внутри docker‑сети)
- **Kafka (Redpanda, Kafka-API совместимый)**: `kafka:9092` (внутри docker‑сети)
- **Redis**: `redis:6379` (внутри docker‑сети)
- **Photo service**: `http://photo-service:8000` (внутри docker‑сети; наружу `http://localhost:8000`)

Перед запуском:

1. Добавьте в `hosts` запись `127.0.0.1 shop.local` (для HTTPS).
2. SSL‑сертификаты берутся из `nginx/ssl` (в репозитории уже есть `localhost.crt/key`).

---

## 🔐 Взаимодействие сервисов

- `hardware_store` (Go) — REST API магазина, интегрирован с `auth-service` по gRPC.
- `store_python` (Flask) — REST API магазина, интегрирован с `auth-service` по gRPC.
- `auth-service` — общий gRPC‑сервис авторизации (JWT, Login, Register, Validate и др.).

**Объединённая авторизация:** оба магазина используют один auth-service. Регистрация в одном магазине даёт доступ к обоим — токен валиден для Go и Python API.
- Контракты в `protos/proto/sso/sso.proto`, сгенерированный код — в `protos/gen/...`.

---

## 🧪 Тестирование и разработка

- Юнит‑ и интеграционные тесты для `auth-service` находятся в `auth-service/tests`.
- Для изменения контрактов:
  1. Меняете `.proto` в `protos/proto/sso/`.
  2. Генерируете код (см. `protos/README.md`).
  3. Обновляете использование в `auth-service` и `hardware_store`.

---

## 📚 Полезные ссылки внутри репозитория

- `hardware_store/README.md` — документация по Go‑магазину.
- `store_python/README.md` — документация по Python‑магазину.
- `protos/README.md` — работа с контрактами и генерацией gRPC‑кода.
- `auth-service/internal` — реализация сервиса авторизации, JWT и gRPC‑сервер.
