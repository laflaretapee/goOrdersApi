# go-orders-api

Небольшой backend-сервис на Go: создание заказа, получение заказа по ID и список заказов.

В проекте есть:
- PostgreSQL в качестве хранилища;
- разделение на HTTP-слой, доменный сервис и слой доступа к БД;
- конфигурация через переменные окружения;
- тесты для бизнес-логики и HTTP-обработчиков.

## Стек

- Go 1.26+
- `net/http`
- PostgreSQL 16
- `database/sql`
- `pgx` driver

## Структура

```text
.
├── cmd/api/main.go
├── internal/config
├── internal/httpapi
├── internal/order
├── internal/storage/postgres
├── migrations
├── docker-compose.yml
└── README.md
```

## Что умеет API

- `POST /orders` создать заказ
- `GET /orders/{id}` получить один заказ
- `GET /orders` получить список заказов

## Модель заказа

```json
{
  "id": 1,
  "customer_name": "Ivan Petrov",
  "item": "iPhone 15",
  "quantity": 1,
  "price_cents": 8999900,
  "status": "new",
  "created_at": "2026-03-27T12:00:00Z"
}
```

`price_cents` хранится в копейках, чтобы не тащить ошибки округления с `float`.

## Быстрый старт

1. Поднять PostgreSQL:

```bash
docker compose up -d db
```

2. Запустить API:

```bash
go run ./cmd/api
```

По умолчанию используется:

```bash
PORT=8080
DATABASE_URL=postgres://orders_user:orders_password@localhost:5432/orders_db?sslmode=disable
```

## Примеры запросов

Создать заказ:

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Ivan Petrov",
    "item": "iPhone 15",
    "quantity": 1,
    "price_cents": 8999900
  }'
```

Получить заказ по ID:

```bash
curl http://localhost:8080/orders/1
```

Получить все заказы:

```bash
curl http://localhost:8080/orders
```

## Тесты

```bash
go test ./...
```
