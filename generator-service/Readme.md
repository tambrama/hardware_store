## Generator Service (Go)

Сервис генерирует события об изменениях товаров и отправляет их в Kafka. 

### Что делает:
- Каждые N секунд (по умолчанию 10 сек) выбирает случайный товар из БД
- Изменяет цену на ±10%
- Изменяет количество на ±15 единиц
- Отправляет события в Kafka топик `product-updates`

### Конфигурация (generator-service/config/local.yaml):
```yaml
database_url: "postgres://postgres:postgres@host.docker.internal:5433/hardwarestore"
kafka_brokers:
  - "host.docker.internal:9093"
kafka_topic: "product-updates"
generator_interval: "10s"       # Интервал между событиями
price_change_percent: 10.0      # Процент изменения цены
stock_change_amount: 15         # Количество единиц изменения
```

### Запуск:
```bash
cd src
docker-compose up generator-service
```
