# Marketplace Microservices

Микросервисная архитектура для платформы электронной коммерции (маркетплейса) на Go.

## Описание проекта

Этот проект представляет собой распределенную систему микросервисов для управления онлайн-маркетплейсом, написанную на Go. Каждый сервис отвечает за определенную бизнес-функцию и может разрабатываться, развертываться и масштабироваться независимо.

## Структура проекта

```
marketplace-microservices/
├── infrastructure/      # Инфраструктурные компоненты (API Gateway, Service Discovery)
├── order-service/       # Сервис управления заказами
├── payment-service/     # Сервис обработки платежей
├── product-service/     # Сервис управления товарами и каталогом
├── proto/              # Protocol Buffers определения для gRPC
├── shared/             # Общие библиотеки и утилиты
├── user-service/       # Сервис управления пользователями и аутентификацией
└── warehouse-service/  # Сервис управления складом и инвентаризацией
```

## Микросервисы

### User Service
Управление пользователями, регистрация, аутентификация и авторизация.

**Основные функции:**
- Регистрация и авторизация пользователей
- Управление профилями
- JWT токены и управление сессиями
- Роли и разграничение доступа

### Product Service
Управление каталогом товаров и их характеристиками.

**Основные функции:**
- CRUD операции для товаров
- Управление категориями
- Поиск и фильтрация товаров
- Управление изображениями товаров

### Order Service
Обработка и управление заказами клиентов.

**Основные функции:**
- Создание и отслеживание заказов
- Управление статусами заказов
- История заказов
- Интеграция с другими сервисами

### Payment Service
Обработка платежей и транзакций.

**Основные функции:**
- Интеграция с платежными системами
- Обработка платежей
- Возвраты и отмены
- Хранение платежной информации

### Warehouse Service
Управление складскими запасами и логистикой.

**Основные функции:**
- Учет товарных остатков
- Резервирование товаров
- Управление поставками
- Отслеживание перемещений товаров

### Infrastructure
Инфраструктурные компоненты для поддержки микросервисов.

**Компоненты:**
- API Gateway (на базе Gin/Echo)
- Service Discovery (Consul)
- Configuration Management
- Health Checks & Metrics

## Технологический стек

- **Language:** Go 1.21+
- **Web Framework:** Gin / Echo / Chi
- **gRPC:** google.golang.org/grpc
- **Database:** PostgreSQL, MongoDB
- **Cache:** Redis
- **Message Broker:** RabbitMQ / Apache Kafka (Sarama)
- **Service Discovery:** Consul / Etcd
- **Containerization:** Docker, Docker Compose
- **Orchestration:** Kubernetes
- **Monitoring:** Prometheus, Grafana
- **Logging:** Zerolog / Zap
- **Tracing:** Jaeger / OpenTelemetry

## Требования

- Go 1.21 или выше
- Docker 20.10+
- Docker Compose 1.29+
- PostgreSQL 14+
- Redis 7+
- RabbitMQ 3.12+ или Kafka 3.x
- Make (опционально)

## Быстрый старт

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd marketplace-microservices
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите инфраструктуру (базы данных, брокер сообщений):
```bash
docker-compose up -d postgres redis rabbitmq
```

4. Запустите все сервисы:
```bash
docker-compose up -d
```

5. Проверьте статус сервисов:
```bash
docker-compose ps
```

6. Доступ к сервисам:
- API Gateway: http://localhost:8080
- Consul UI: http://localhost:8500
- RabbitMQ Management: http://localhost:15672 (guest/guest)
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

## Разработка

### Запуск отдельного сервиса локально

```bash
cd user-service
go run cmd/main.go
```

### Сборка сервиса

```bash
cd user-service
go build -o bin/user-service cmd/main.go
```

### Запуск тестов

```bash
# Все тесты
go test ./...

# С покрытием
go test -cover ./...

# Конкретный пакет
go test ./internal/handler
```

### Генерация gRPC кода из proto файлов

```bash
# Установка protoc и плагинов
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация
cd proto
protoc --go_out=. --go-grpc_out=. *.proto
```

### Линтинг и форматирование

```bash
# Форматирование кода
go fmt ./...

# Линтинг
golangci-lint run

# Установка golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Структура сервиса (пример)

```
user-service/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── config/              # Конфигурация
│   ├── domain/              # Доменные модели
│   ├── handler/             # HTTP/gRPC handlers
│   ├── repository/          # Слой доступа к данным
│   ├── service/             # Бизнес-логика
│   └── middleware/          # Middleware
├── pkg/
│   └── utils/               # Утилиты
├── migrations/              # SQL миграции
├── Dockerfile
├── go.mod
└── go.sum
```

## API документация

### REST API
API документация генерируется с помощью Swagger/OpenAPI:
- Swagger UI: http://localhost:8080/swagger/index.html

Для генерации документации используйте:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go
```

### gRPC API
Документация Protocol Buffers находится в папке `/proto`. Для тестирования gRPC API используйте:
- grpcurl
- BloomRPC
- Postman (с поддержкой gRPC)

## Конфигурация

Каждый сервис использует переменные окружения или конфигурационные файлы:

```env
# .env пример
SERVER_PORT=8081
DATABASE_URL=postgres://user:password@localhost:5432/dbname
REDIS_URL=redis://localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
JWT_SECRET=your-secret-key
LOG_LEVEL=debug
```

Конфигурация может загружаться с помощью:
- viper
- envconfig
- Consul KV Store

## Взаимодействие между сервисами

### Синхронное взаимодействие
- REST API (HTTP/JSON)
- gRPC (Protocol Buffers)

### Асинхронное взаимодействие
- RabbitMQ для event-driven архитектуры
- Kafka для обработки потоков событий

### Паттерны
- Event Sourcing
- CQRS
- Saga Pattern для распределенных транзакций

## Миграции базы данных

Используется migrate или goose:

```bash
# Установка migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Применение миграций
migrate -path ./migrations -database "postgres://localhost/dbname?sslmode=disable" up

# Откат
migrate -path ./migrations -database "postgres://localhost/dbname?sslmode=disable" down 1
```

## Мониторинг и логирование

### Метрики (Prometheus)
Каждый сервис экспортирует метрики:
- HTTP запросы (latency, status codes)
- gRPC вызовы
- Database connections
- Custom business metrics

Эндпоинт метрик: http://localhost:8081/metrics

### Логирование
Структурированное логирование с помощью zerolog или zap:
```go
log.Info().
    Str("service", "user-service").
    Int("user_id", 123).
    Msg("User created successfully")
```

### Трейсинг
Distributed tracing с OpenTelemetry и Jaeger:
- Jaeger UI: http://localhost:16686

## Docker

### Сборка образа

```bash
docker build -t user-service:latest ./user-service
```

### Запуск контейнера

```bash
docker run -p 8081:8081 \
  -e DATABASE_URL=postgres://... \
  user-service:latest
```


