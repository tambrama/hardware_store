# Store Python (Flask)

REST API интернет-магазина бытовой техники на Python (Flask, PostgreSQL, Flasgger/Swagger).

## Структура

- `src/app/` — приложение Flask
- `src/seed.py` — заполнение БД тестовыми данными
- `run.py` — точка входа (запуск из корня store_python)
- `scripts/run_replica.ps1` — запуск реплики для локальной балансировки

Доп. сервисы (через общий `src/docker-compose.yml`):
- `photo-service` — отдельный сервис хранения фото (FastAPI) + шардирование по UUID на 4 базы Postgres.

## Запуск (локально)

```bash
cd store_python
pip install -r src/requirements.txt
export DATABASE_URI=postgresql://postgres:postgres@localhost:5432/shopapi
export PHOTO_SERVICE_URL=http://localhost:8000
python run.py
```

## Seed (тестовые данные)

```bash
cd store_python
$env:PYTHONPATH="src"; python src/seed.py
# с перезаполнением: python src/seed.py --force
```

## Docker

Сборка через `docker-compose` в корне `src/` — сервисы `backend`, `backend-replica1`, `backend-replica2`, `photo-service`, `db-photos`.

> Если раньше поднимали проект со старой схемой (фото в shopapi БД), проще всего сделать `docker compose down -v` (или `make docker-rebuild`, если у вас есть) — т.к. менялась модель `product.image_id` (убран FK на локальную таблицу images).
