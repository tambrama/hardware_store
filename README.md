<h1>🏪 Hardware Store API </h1>

RESTful API для управления интернет-магазином строительных материалов.

---
<!-- omit in toc -->
<h2>📖 Содержание</h2>

- [🌟 Быстрый старт](#-быстрый-старт)
- [🚀 Технологии](#-технологии)
- [📁 Структура проекта](#-структура-проекта)
- [🚀 Установка и запуск](#-установка-и-запуск)
  - [🛠️ C Docker Compose (рекомендуется)](#️-c-docker-compose-рекомендуется)
  - [🚦 Остановка](#-остановка)
  - [🔄 Полный перезапуск (сброс БД)](#-полный-перезапуск-сброс-бд)
- [🏗️ Архитектура](#️-архитектура)
  - [🔄 Слои:](#-слои)
- [📚 API Документация](#-api-документация)
- [📝 Makefile команды](#-makefile-команды)

---
## 🌟 Быстрый старт
```
# Запуск с Docker Compose (рекомендуется)
docker-compose up --build

# Сервер будет доступен по адресу:
# http://localhost:8081
```
---
## 🚀 Технологии

- **Go 1.25+** - язык программирования
- **Gin** - веб-фреймворк
- **PostgreSQL** - реляционная база данных
- **Docker & Docker Compose** - контейнеризация и оркестрация
- **Swaggo** - автоматическая генерация Swagger документации
- **Validator v10** - валидация входных данных
- **UUID** - генерация уникальных идентификаторов
- **Uber FX** - dependency injection (через кастомный DI)
- **Многослойная архитектура** - архитектура с разделением на слои
---
## 📁 Структура проекта
```
│   .env
│   docker-compose.yml
│   go.mod
│   go.sum
│   Makefile
│   README.md
│
├───build
│   ├───migrator
│   │       Dockerfile
│   │
│   └───server
│           Dockerfile
│
├───cmd
│   └───app
│           main.go
│
├───docs
│       docs.go
│       swagger.json
│       swagger.yaml
│
├───internal
│   ├───app
│   │       app.go
│   │
│   ├───config
│   │       config.go
│   │       local.yaml
│   │
│   ├───di
│   │       di.go
│   │
│   ├───logger
│   │       logger.go
│   │
│   ├───model
│   │   ├───address
│   │   │       address.go
│   │   │
│   │   ├───category
│   │   │       category.go
│   │   │
│   │   ├───client
│   │   │       client.go
│   │   │
│   │   ├───images
│   │   │       images.go
│   │   │
│   │   ├───product
│   │   │       product.go
│   │   │
│   │   ├───supplier
│   │   │       supplier.go
│   │   │
│   │   └───tx
│   │           manager.go
│   │
│   ├───server
│   │       server.go
│   │
│   ├───service
│   │   ├───address_service
│   │   │       service.go
│   │   │       service_impl.go
│   │   │
│   │   ├───category_service
│   │   │       service.go
│   │   │       service_impl.go
│   │   │
│   │   ├───client_service
│   │   │       service.go
│   │   │       service_impl.go
│   │   │
│   │   ├───images_service
│   │   │       service.go
│   │   │       service_impl.go
│   │   │
│   │   ├───product_service
│   │   │       service.go
│   │   │       service_impl.go
│   │   │
│   │   └───supplier_service
│   │           service.go
│   │           service_impl.go
│   │
│   ├───storage
│   │   │   storage.go
│   │   │
│   │   └───postgres
│   │       │   address.go
│   │       │   category.go
│   │       │   client.go
│   │       │   images.go
│   │       │   postgres.go
│   │       │   product.go
│   │       │   supplier.go
│   │       │   tx.go
│   │       │
│   │       ├───dto
│   │       │       dto.go
│   │       │
│   │       └───mapper
│   │               address.go
│   │               category.go
│   │               client.go
│   │               images.go
│   │               product.go
│   │               supplier.go
│   │
│   └───web
│       │   router.go
│       │
│       ├───dto
│       │       dto.go
│       │       error.go
│       │
│       ├───handler
│       │       address_handler.go
│       │       category_handler.go
│       │       client_handler.go
│       │       image_handler.go
│       │       product_handler.go
│       │       supplier_handler.go
│       │
│       └───mapper
│               mapper.go
│
└───migrations
        000001_create_table.sql
        000002_seed_test_data.sql
```
---

## 🚀 Установка и запуск

### 🛠️ C Docker Compose (рекомендуется)
1. **Создайте файл `.env` в корне проекта:**
```
CONFIG_PATH=/app/internal/config/local.yaml
# Порт HTTP-сервера
SERVER_PORT=8081

# Строка подключения к PostgreSQL (для локального запуска)
APP_ENV=local
DATABASE_URL=postgres://postgres:password@db:5432/hardwarestore?sslmode=disable
PSQL_NAME=hardwarestore
PSQL_USER=postgres
PSQL_PASSWORD=password
```
2. **Запустите приложение:**
```
docker-compose up --build
```
3. **Приложение будет доступно:**
    → Сервер: [http://localhost:8081](http://localhost:8081/)
    → База данных: PostgreSQL на порту 5432

###  🚦 Остановка
```
docker-compose down
```
### 🔄 Полный перезапуск (сброс БД)
```
make -C src docker-rebuild
```

> ⚠️ Все данные будут удалены!

---

##  🏗️ Архитектура
Проект использует многослойную архитектуру с чётким разделением ответственности:
```
┌─────────────────────────────────────────────┐
│              HTTP Handler (Gin)             │
├─────────────────────────────────────────────┤
│              Service Layer                  │
│  (бизнес-логика, валидация, транзакции)     │
├─────────────────────────────────────────────┤
│            Repository Layer                 │
│       (работа с PostgreSQL)                 │
├─────────────────────────────────────────────┤
│              Domain Models                  │
│         (чистые бизнес-объекты)             │
└─────────────────────────────────────────────┘
```
### 🔄 Слои:
Handler - обработка HTTP запросов, валидация входных данных
Service - бизнес-логика, координация операций, транзакции
Storage/Repository - абстракция работы с базой данных
Model - доменные сущности
DTO/Mapper - преобразование данных между слоями

---

## 📚 API Документация
После запуска приложения документация доступна по адресу:
```
http://localhost:8081/swagger/index.html
```
Для генерации Swagger документации:
```
# Установить swag (если не установлен)
go install github.com/swaggo/swag/cmd/swag@latest

# Сгенерировать документацию
swag init

# Или через Makefile
make docs
```

---

## 📝 Makefile команды
```
# Основные команды:
make run         # Запуск сервера
make build       # Сборка бинарника
make clean       # Очистка собранных файлов

# Docker команды:
make docker-up   # Запуск через Docker Compose
make docker-down # Остановка Docker контейнеров
make docker-rebuild # Полная пересборка

# Утилиты:
make fmt         # Форматирование кода
make update_mod  # Обновление зависимостей
```
---