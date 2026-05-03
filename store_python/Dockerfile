# Бэкенд ShopAPI (Flask)
# Контекст сборки: src/ (родитель store_python)
FROM python:3.12-slim

WORKDIR /app

# Системные зависимости для psycopg2 (PostgreSQL)
RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential libpq-dev && \
    rm -rf /var/lib/apt/lists/*

COPY store_python/src/requirements.txt ./requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Protos для gRPC auth-client
COPY protos/gen/python /app/protos_python

COPY store_python/src/ ./src/
COPY store_python/run.py .

ENV PYTHONPATH=/app/src
ENV PROTOS_PATH=/app/protos_python
ENV FLASK_APP=app:create_app

EXPOSE 5000

CMD ["python", "run.py"]
