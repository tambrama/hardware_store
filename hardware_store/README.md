<h1>🏪 Hardware Store API </h1>

RESTful API для управления интернет-магазином бытовой техники.

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
- [🌐 Nginx конфигурация](#-nginx-конфигурация)
    - [📋 Функционал Nginx](#-функционал-nginx)
    - [🔒 HTTPS и SSL](#-https-и-ssl)
    - [⚖️ Балансировка нагрузки](#️-балансировка-нагрузки)
    - [🗄️ Кеширование](#️-кеширование)
    - [📦 Gzip сжатие](#-gzip-сжатие)
    - [🔄 Reverse Proxy](#-reverse-proxy)
    - [📊 Мониторинг](#-мониторинг)
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
# https://local.hardwarestore.com/
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
├───nginx
│   ├───ssl
│   │       local.hardwarestore.com.crt     
│   │       local.hardwarestore.com.key     
│   │
│   └───nginx.conf                          
│
├───static
│       index.html                          
│       image.png                           
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

# pgAdmin
PGADMIN_EMAIL=admin@test.com
PGADMIN_PASSWORD=mypassword123

# Domain
DOMAIN_NAME=local.hardwarestore.com
```
2. **Добавьте запись в hosts файл:**
```
# Windows (C:\Windows\System32\drivers\etc\hosts)
127.0.0.1 local.hardwarestore.com

# Linux/Mac (/etc/hosts)
127.0.0.1 local.hardwarestore.com
```
3**Создайте самоподписанный SSL сертификат:**
```
cd nginx/ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout local.hardwarestore.com.key \
  -out local.hardwarestore.com.crt \
  -subj "/CN=local.hardwarestore.com"
```
4**Запустите приложение:**
```
docker-compose up --build
```
3. **Приложение будет доступно:**
   → Основной API: https://local.hardwarestore.com/api/...
   → Swagger UI: https://local.hardwarestore.com/api/v1/
   → pgAdmin: https://local.hardwarestore.com/admin/
   → Статика: https://local.hardwarestore.com/
   → Статус Nginx: https://local.hardwarestore.com/status
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

## 🌐 Nginx конфигурация
### 📋 Функционал Nginx
| Функция | Описание |
|---------|----------|
| Reverse Proxy | Проксирование запросов к бэкенд серверам |
| Балансировка нагрузки | GET запросы распределяются 2:1:1 между 3 бэкендами |
| SSL Termination | HTTPS (порт 443) с самоподписанным сертификатом |
| Кеширование | Статические файлы кешируются на 30 дней |
| Gzip сжатие | Текстовые ресурсы сжимаются, изображения - нет |
| Раздача статики | index.html и image.png по пути `/` |
| Мониторинг | Статистика Nginx по пути `/status` |
### 🔒 HTTPS и SSL
* Самоподписанный SSL сертификат для local.hardwarestore.com
* HTTP (порт 80) и HTTPS (порт 443) в одном server блоке
* Поддержка HTTP/2 для улучшения производительности
* Современные протоколы TLSv1.2 и TLSv1.3

### ⚖️ Балансировка нагрузки
Запущено 3 инстанса бэкенда:
* app-main (вес 2) - основной сервер, полный доступ к БД
* app-read-1 (вес 1) - read-only доступ к БД
* app-read-2 (вес 1) - read-only доступ к БД
  Логика балансировки:
* GET запросы к /api/v1/* и /api/v2/* распределяются в соотношении 2:1:1
* Все остальные методы (POST, PUT, DELETE) идут только на основной сервер
* Запросы к /api/* автоматически преобразуются в /api/v1/*
### 🗄️ Кеширование
* Кеширование статических файлов (30 дней)
* Заголовки Cache-Control: public, immutable
* API запросы не кешируются
* Отдельное кеширование для Swagger UI (5 минут)

### 📦 Gzip сжатие
* Сжатие текстовых форматов: HTML, CSS, JS, JSON, XML
* Изображения (PNG, JPEG, GIF и др.) НЕ сжимаются
* Минимальный размер для сжатия: 1024 байта
* Уровень сжатия: 6 (оптимальный баланс)

### 🔄 Reverse Proxy
```
Пользователь → [HTTPS] → Nginx (443) → [HTTP] → Бэкенды (8081)
↓
🔄 Балансировка
↓
┌──────────────┼──────────────┐
app-main        app-read-1     app-read-2
(вес 2)          (вес 1)        (вес 1)
```
### 📊 Мониторинг
Страница /status показывает статистику Nginx:
* Активные соединения
* Всего принятых/обработанных запросов
* Текущие чтение/запись/ожидание

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