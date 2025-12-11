# Примеры использования API сервиса сокращения URL

Этот пакет содержит примеры работы со всеми основными эндпоинтами сервиса сокращения URL.

## Запуск примеров

Для просмотра примеров используйте команду `go doc`:

```bash
# Просмотр всех примеров
go doc -all examples

# Просмотр конкретного примера
go doc examples.Example_shortenURL
```

## Список примеров

### Example_shortenURL
Демонстрирует создание короткого URL через текстовый запрос (Content-Type: text/plain).

**Эндпоинт:** `POST /`

**Пример использования:**
```go
go doc examples.Example_shortenURL
```

### Example_shortenURLJSON
Демонстрирует создание короткого URL через JSON запрос.

**Эндпоинт:** `POST /api/shorten`

**Пример запроса:**
```json
{
  "url": "https://practicum.yandex.ru"
}
```

### Example_shortenBatch
Демонстрирует пакетное создание коротких URL.

**Эндпоинт:** `POST /api/shorten/batch`

**Пример запроса:**
```json
[
  {
    "correlation_id": "req1",
    "original_url": "https://practicum.yandex.ru"
  },
  {
    "correlation_id": "req2",
    "original_url": "https://yandex.ru"
  }
]
```

### Example_getURL
Демонстрирует получение оригинального URL по короткому с редиректом.

**Эндпоинт:** `GET /:id`

### Example_getUserURLs
Демонстрирует получение всех URL пользователя.

**Эндпоинт:** `GET /api/user/urls`

**Пример ответа:**
```json
[
  {
    "short_url": "http://localhost:8080/abc123",
    "original_url": "https://practicum.yandex.ru"
  }
]
```

### Example_deleteUserURLs
Демонстрирует удаление (мягкое) URL пользователя.

**Эндпоинт:** `DELETE /api/user/urls`

**Пример запроса:**
```json
["abc123", "def456"]
```

### Example_ping
Демонстрирует проверку доступности базы данных (healthcheck).

**Эндпоинт:** `GET /ping`

## Аутентификация

Все примеры (кроме `Example_ping` и `Example_getURL`) требуют аутентификации через cookie `user_id` с подписью. В примерах используется `AuthService` для создания подписанных cookies.

## Дополнительная документация

Для просмотра полной документации по API используйте:

```bash
# Документация по хендлерам
go doc -all internal/handler

# Документация по сервису
go doc -all internal/service

# Документация по моделям
go doc -all internal/model
```
