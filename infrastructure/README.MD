# Kafka Docker Setup

Локальная установка Apache Kafka с Zookeeper и Kafka UI с использованием Docker Compose.

## 📋 Содержание

- [Требования](#требования)
- [Структура проекта](#структура-проекта)
- [Компоненты](#компоненты)
- [Быстрый старт](#быстрый-старт)
- [Конфигурация](#конфигурация)
- [Использование](#использование)
- [Управление](#управление)
- [Порты](#порты)
- [Troubleshooting](#troubleshooting)

## 🔧 Требования

- Docker 20.10+
- Docker Compose 1.29+ (или Docker Compose V2)
- Минимум 4GB RAM
- Свободные порты: 2181, 9092, 29092, 8080

## 📁 Структура проекта

```
.
├── docker-compose.yml    # Основной файл с настройками запуска
├── config.yaml          # Базовая конфигурация сервисов
└── README.md           # Документация
```

## 🐳 Компоненты

### Zookeeper
- **Версия**: 7.5.0
- **Порт**: 2181
- **Описание**: Координатор для Kafka кластера

### Kafka
- **Версия**: 7.5.0
- **Порты**:
    - 9092 (внешний доступ)
    - 29092 (внутрикластерный)
- **Описание**: Брокер сообщений Apache Kafka

### Kafka UI
- **Версия**: latest
- **Порт**: 8080
- **Описание**: Web-интерфейс для управления Kafka

## 🚀 Быстрый старт

### 1. Клонирование и запуск

```bash
# Запуск всех сервисов
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f
```

### 2. Проверка работоспособности

Откройте в браузере: http://localhost:8080

Вы должны увидеть Kafka UI с подключенным кластером `local`.

## ⚙️ Конфигурация

### config.yaml

Содержит базовые настройки сервисов (образы, переменные окружения, volumes). Используется как основа для наследования.

### docker-compose.yml

Содержит специфичные настройки для запуска:
- Имена контейнеров
- Порты
- Зависимости между сервисами
- Наследует конфигурацию из `config.yaml`

### Переменные окружения Kafka

| Переменная | Значение | Описание |
|------------|----------|----------|
| `KAFKA_BROKER_ID` | 1 | ID брокера |
| `KAFKA_ZOOKEEPER_CONNECT` | zookeeper:2181 | Подключение к Zookeeper |
| `KAFKA_ADVERTISED_LISTENERS` | PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092 | Адреса для подключения |
| `KAFKA_AUTO_CREATE_TOPICS_ENABLE` | true | Автосоздание топиков |

## 💡 Использование

### Создание топика

```bash
docker exec -it kafka kafka-topics --create \
  --topic my-topic \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1
```

### Список топиков

```bash
docker exec -it kafka kafka-topics --list \
  --bootstrap-server localhost:9092
```

### Отправка сообщения

```bash
docker exec -it kafka kafka-console-producer \
  --topic my-topic \
  --bootstrap-server localhost:9092
```

### Чтение сообщений

```bash
docker exec -it kafka kafka-console-consumer \
  --topic my-topic \
  --from-beginning \
  --bootstrap-server localhost:9092
```

### Подключение из приложения

```python
# Python (kafka-python)
from kafka import KafkaProducer

producer = KafkaProducer(bootstrap_servers='localhost:9092')
producer.send('my-topic', b'Hello Kafka!')
```

```javascript
// Node.js (kafkajs)
const { Kafka } = require('kafkajs');

const kafka = new Kafka({
  clientId: 'my-app',
  brokers: ['localhost:9092']
});
```

```java
// Java
Properties props = new Properties();
props.put("bootstrap.servers", "localhost:9092");
props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");

KafkaProducer<String, String> producer = new KafkaProducer<>(props);
```

## 🎮 Управление

### Запуск

```bash
# Запуск в фоновом режиме
docker-compose up -d

# Запуск с просмотром логов
docker-compose up
```

### Остановка

```bash
# Остановка всех сервисов
docker-compose down

# Остановка с удалением volumes (УДАЛИТ ВСЕ ДАННЫЕ!)
docker-compose down -v
```

### Перезапуск

```bash
# Перезапуск всех сервисов
docker-compose restart

# Перезапуск конкретного сервиса
docker-compose restart kafka
```

### Логи

```bash
# Все логи
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f kafka

# Последние 100 строк
docker-compose logs --tail=100 kafka
```

### Масштабирование

```bash
# Запуск нескольких брокеров Kafka (требует изменения конфигурации)
docker-compose up -d --scale kafka=3
```

## 🔌 Порты

| Сервис | Порт | Назначение |
|--------|------|------------|
| Zookeeper | 2181 | Client port |
| Kafka | 9092 | Внешний доступ (PLAINTEXT_HOST) |
| Kafka | 29092 | Внутренний доступ (PLAINTEXT) |
| Kafka UI | 8080 | Web интерфейс |

## 🔍 Troubleshooting

### Kafka не запускается

```bash
# Проверьте логи
docker-compose logs kafka

# Проверьте, запущен ли Zookeeper
docker-compose ps zookeeper

# Перезапустите сервисы
docker-compose restart
```

### Ошибка "Address already in use"

Порты уже заняты. Остановите процессы или измените порты в `docker-compose.yml`:

```yaml
ports:
  - "9093:9092"  # Измените внешний порт
```

### Недостаточно памяти

Увеличьте память для Docker (Docker Desktop → Settings → Resources → Memory).

### Kafka UI не показывает данные

```bash
# Проверьте подключение
docker exec -it kafka-ui env | grep KAFKA

# Перезапустите Kafka UI
docker-compose restart kafka-ui
```

### Очистка всех данных

```bash
# Остановка и удаление всех данных
docker-compose down -v

# Удаление образов (опционально)
docker-compose down --rmi all

# Новый запуск
docker-compose up -d
```