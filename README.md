# Microservices — Космические запчасти

Учебный проект из курса **«Микросервисы, как в BigTech 2.0»**.  
Реализует упрощённую e-commerce платформу для продажи деталей космических кораблей.

Go Workspace объединяет четыре независимых модуля: `shared`, `inventory`, `payment`, `order`.

---

## Архитектура

```
┌─────────────────────────────────────────────────────────┐
│                        Клиент                           │
│              HTTP REST (OpenAPI 3.0)                    │
└─────────────────────────┬───────────────────────────────┘
                          │ :8080
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   Order Service                         │
│         HTTP/REST — chi router + ogen (gen.)            │
│                                                         │
│  POST   /api/v1/orders          — создать заказ         │
│  GET    /api/v1/orders/{uuid}   — получить заказ        │
│  POST   /api/v1/orders/{uuid}/pay    — оплатить         │
│  POST   /api/v1/orders/{uuid}/cancel — отменить         │
│                                                         │
│  in-memory хранилище  (map + sync.RWMutex)              │
└──────────┬──────────────────────────┬───────────────────┘
           │ gRPC :50051              │ gRPC :50052
           ▼                          ▼
┌──────────────────────┐  ┌──────────────────────────────┐
│  Inventory Service   │  │      Payment Service         │
│  gRPC сервер         │  │      gRPC сервер             │
│  :50051              │  │      :50052                  │
│                      │  │                              │
│  GetPart(uuid)       │  │  PayOrder(order, user,       │
│  ListParts(filter)   │  │           method)            │
│                      │  │  → transaction_uuid          │
│  in-memory каталог   │  │  (заглушка, всегда успех)    │
│  8 тестовых деталей  │  │                              │
└──────────────────────┘  └──────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                  shared (общий модуль)                  │
│                                                         │
│  pkg/api/order/v1/       — сгенерированный ogen-код     │
│  pkg/proto/inventory/v1/ — сгенерированный protobuf     │
│  pkg/proto/payment/v1/   — сгенерированный protobuf     │
│                                                         │
│  api/order/v1/           — OpenAPI YAML-спецификация    │
│  proto/inventory/v1/     — .proto файл Inventory        │
│  proto/payment/v1/       — .proto файл Payment          │
└─────────────────────────────────────────────────────────┘
```

---

## Структура репозитория

```
.
├── go.work                   # Go Workspace — объединяет все модули
├── Taskfile.yml              # Задачи: генерация, линтинг, форматирование
├── .golangci.yml             # Конфигурация golangci-lint v2
├── package.json              # @redocly/cli — линтинг OpenAPI схем
│
├── shared/                   # Общий модуль: контракты и сгенерированный код
│   ├── proto/                # Исходные .proto файлы
│   │   ├── inventory/v1/inventory.proto
│   │   └── payment/v1/payment.proto
│   ├── api/                  # OpenAPI спецификации (разбиты на компоненты)
│   │   └── order/v1/         # order.openapi.yaml + paths/ + components/
│   └── pkg/                  # Сгенерированный Go-код (не редактировать вручную)
│       ├── api/order/v1/     # ogen-генерация HTTP-сервера/клиента заказов
│       └── proto/            # protoc-генерация gRPC клиентов/серверов
│           ├── inventory/v1/
│           └── payment/v1/
│
├── inventory/                # Сервис складского учёта (gRPC :50051)
│   └── cmd/main.go
│
├── payment/                  # Сервис обработки платежей (gRPC :50052)
│   └── cmd/main.go
│
└── order/                    # Сервис заказов (HTTP REST :8080)
    └── cmd/main.go
```

---

## Описание модулей

### `shared` — общий контрактный модуль

Центральный модуль, от которого зависят все остальные. Содержит только контракты и сгенерированный по ним код. Не содержит бизнес-логики.

| Путь | Описание |
|------|----------|
| `proto/inventory/v1/inventory.proto` | gRPC-контракт сервиса склада: `GetPart`, `ListParts`, фильтрация по UUID/имени/категории/стране/тегам |
| `proto/payment/v1/payment.proto` | gRPC-контракт платёжного сервиса: `PayOrder` → `transaction_uuid` |
| `api/order/v1/` | OpenAPI 3.0 спецификация REST API заказов, разбитая на компоненты (paths, schemas, errors) |
| `pkg/proto/*/` | Автогенерация из `.proto` через `buf` + `protoc-gen-go` + `protoc-gen-go-grpc` |
| `pkg/api/order/v1/` | Автогенерация HTTP-сервера и клиента из OpenAPI через `ogen` |

**Категории деталей (Inventory):** `ENGINE`, `FUEL`, `PORTHOLE`, `WING`  
**Методы оплаты (Payment):** `CARD`, `SBP`, `CREDIT_CARD`, `INVESTOR_MONEY`

---

### `inventory` — сервис складского учёта

**Протокол:** gRPC  
**Порт:** `:50051`

Хранит каталог деталей для космических кораблей. Данные живут в памяти (`map[string]*Part` + `sync.RWMutex`) — при рестарте инициализируются тестовым набором из 8 деталей.

**gRPC методы:**
- `GetPart(uuid)` — получить одну деталь по UUID
- `ListParts(filter)` — получить список деталей с фильтрацией

**Фильтрация деталей** поддерживает одновременно несколько критериев (UUID, имя, категория, страна производителя, теги). Фильтры применяются через логику AND, внутри каждого фильтра — OR.

**Тестовые данные включают:**

| Название | Категория | Цена |
|----------|-----------|------|
| Ионный двигатель X-2000 | ENGINE | 150 000 ₽ |
| Плазменный двигатель P-500 | ENGINE | 200 000 ₽ |
| Криогенное топливо H2-O2 | FUEL | 50 000 ₽ |
| Ядерное топливо U-235 | FUEL | 300 000 ₽ |
| Кварцевое окно QW-100 | PORTHOLE | 25 000 ₽ |
| Бронированное окно BW-200 | PORTHOLE | 40 000 ₽ |
| Солнечная панель SP-500 | WING | 75 000 ₽ |
| Аэродинамическое крыло AW-300 | WING | 60 000 ₽ |

Включена gRPC Reflection для отладки через `grpcurl`.

---

### `payment` — сервис обработки платежей

**Протокол:** gRPC  
**Порт:** `:50052`

Минималистичная заглушка платёжного провайдера. При любом запросе `PayOrder` генерирует UUID транзакции и возвращает успех. Реальной интеграции с платёжными системами нет — намеренно упрощён для учебных целей.

**gRPC методы:**
- `PayOrder(order_uuid, user_uuid, payment_method)` → `transaction_uuid`

Включена gRPC Reflection.

---

### `order` — сервис заказов

**Протокол:** HTTP REST (OpenAPI 3.0)  
**Порт:** `:8080`  
**Роутер:** `go-chi/chi v5`  
**Генератор кода:** `ogen`

Центральный оркестратор. Принимает REST-запросы, координирует работу `inventory` и `payment` через gRPC. Хранит заказы в памяти.

**REST эндпоинты:**

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/api/v1/orders` | Создать заказ: запрашивает цены деталей у Inventory, считает `total_price`, сохраняет со статусом `PENDING_PAYMENT` |
| `GET` | `/api/v1/orders/{order_uuid}` | Получить заказ по UUID |
| `POST` | `/api/v1/orders/{order_uuid}/pay` | Оплатить заказ: проверяет статус, вызывает Payment Service, переводит в `PAID` |
| `POST` | `/api/v1/orders/{order_uuid}/cancel` | Отменить заказ: переводит в `CANCELLED` (нельзя отменить уже оплаченный) |

**Жизненный цикл заказа:**

```
PENDING_PAYMENT ──► PAID
PENDING_PAYMENT ──► CANCELLED
```

**Graceful shutdown:** при получении `SIGINT`/`SIGTERM` даёт серверу 10 секунд на завершение активных запросов.

**Middleware:** Logger, Recoverer, Timeout (10 сек).

---

## Межсервисное взаимодействие

```
Клиент
  │
  │ POST /api/v1/orders  { user_uuid, part_uuids[] }
  ▼
Order Service
  │
  │ для каждого part_uuid:
  │   gRPC ListParts({ uuids: [part_uuid] })  ──►  Inventory Service
  │                                           ◄──  { parts: [{ price, ... }] }
  │ суммирует цены → total_price
  │ создаёт заказ (PENDING_PAYMENT)
  │
  │ POST /api/v1/orders/{uuid}/pay  { payment_method }
  ▼
Order Service
  │
  │ gRPC PayOrder({ order_uuid, user_uuid, payment_method })  ──►  Payment Service
  │                                                           ◄──  { transaction_uuid }
  │ статус → PAID, сохраняет transaction_uuid
```

---

## Запуск

Запустить каждый сервис в отдельном терминале:

```bash
go run inventory/cmd/main.go   # gRPC :50051
go run payment/cmd/main.go     # gRPC :50052
go run order/cmd/main.go       # HTTP :8080
```

### Требования к инструментам

```bash
brew install go-task   # Taskfile CLI
```

### Основные задачи Taskfile

```bash
task proto:gen    # Регенерировать gRPC код из .proto (через buf)
task ogen:gen      # Регенерировать HTTP код из OpenAPI (через ogen)
task lint              # Запустить golangci-lint
task format               # Форматирование кода (gofumpt + gci)
```

---

## Стек технологий

| Категория | Инструмент | Назначение |
|-----------|------------|------------|
| Язык | Go 1.25.5 | — |
| HTTP роутер | `go-chi/chi v5` | REST сервер в Order |
| HTTP кодогенерация | `ogen v1.12` | Генерация сервера/клиента из OpenAPI |
| gRPC | `google.golang.org/grpc` | Межсервисное взаимодействие |
| Protobuf кодогенерация | `buf` + `protoc-gen-go` | Генерация Go-кода из `.proto` |
| Линтер | `golangci-lint v2` | Статический анализ кода |
| Форматтеры | `gofumpt` + `gci` | Форматирование кода и импортов |
| OpenAPI линтер | `@redocly/cli` | Валидация OpenAPI схем |
| Сборщик задач | `Taskfile` | Автоматизация dev-задач |
