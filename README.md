# Tasks API

Простой REST API для управления задачами (To‑Do list) на Go.

## Возможности

- CRUD операции над задачами:
  - `GET /tasks` — список всех задач
  - `POST /tasks` — создание новой задачи
  - `GET /tasks/{id}` — получение задачи по ID
  - `PUT /tasks/{id}` — обновление задачи
  - `DELETE /tasks/{id}` — удаление задачи
- In‑memory хранилище с потокобезопасностью (`sync.RWMutex`)
- Валидация обязательных полей (например, `title`)
- JSON ответы с кодами статусов:
  - `200 OK` — успешное чтение/обновление
  - `201 Created` — успешное создание
  - `204 No Content` — успешное удаление
  - `400 Bad Request` — неверный формат данных или валидация
  - `404 Not Found` — задача не найдена
  - `405 Method Not Allowed` — метод не поддерживается

## Структура проекта

```text
.
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── handlers
│   │   └── task.go
│   ├── http
│   │   └── middleware.go(не реализовал, но принцип понял, что нужно перенести логику из /internal/handlers/task.go)
│   ├── models
│   │   └── task.go
│   └── storage
│       └── storage.go
├── go.mod
└── README.md
```

## Модель задачи

```json
{
  "id": 1,
  "title": "Buy milk",
  "done": false,
  "created_at": "2026-08-17T13:00:00+03:00"
}
```

Поля:

- `id` (int) — уникальный идентификатор (генерируется сервером)
- `title` (string) — заголовок задачи (обязательное поле)
- `done` (bool) — статус выполнения
- `created_at` (string, RFC3339) — время создания (генерируется сервером)

## Запуск сервера

```bash
go run ./cmd/server
```

Сервер слушает: `http://localhost:8081`

## Тестирование API

### PowerShell (Windows)

#### GET — список всех задач

powershell
```powershell
Invoke-RestMethod -Uri http://localhost:8081/tasks -Method GET
```
curl
```curl
curl http://localhost:8081/tasks
```
браузер
```в браузере
http://localhost:8081/tasks
```

#### POST — создание задачи
powershell
```powershell
$body = @{
    title = "Buy milk"
    done = $false
} | ConvertTo-Json

$r = Invoke-WebRequest -Uri http://localhost:8081/tasks -Method POST -ContentType "application/json" -Body $body -UseBasicParsing
$r.Content
```
curl
```curl
curl -X POST http://localhost:8081/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk", "done": false}'
```
#### GET — получение задачи по ID
powershell
```powershell
Invoke-RestMethod -Uri http://localhost:8081/tasks/1 -Method GET
```
curl
```curl
curl http://localhost:8081/tasks/1
```

#### PUT — обновление задачи
powershell
```powershell
$body = @{
    title = "Buy coffee"
    done = $false
} | ConvertTo-Json

$r = Invoke-WebRequest -Uri http://localhost:8081/tasks/1 -Method PUT -ContentType "application/json" -Body $body -UseBasicParsing
$r.Content
```
curl
```curl
curl -X PUT http://localhost:8081/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy coffee", "done": false}'
```

#### DELETE — удаление задачи
powershell
```powershell
$r = Invoke-WebRequest -Uri http://localhost:8081/tasks/1 -Method DELETE -UseBasicParsing
$r.StatusCode
```
curl
```curl
curl -X DELETE http://localhost:8081/tasks/1 -v
```
---
## Примеры ответов

### Успешное создание (POST)

Статус: `201 Created`

```json
{
  "id": 1,
  "title": "Buy milk",
  "done": false,
  "created_at": "2026-08-17T13:00:00+03:00"
}
```

### Ошибка валидации (POST)

Статус: `400 Bad Request`

```json
{
  "error": "Вы не заполнили заголовок task"
}
```

### Задача не найдена (GET/PUT/DELETE)

Статус: `404 Not Found`

```json
{
  "error": "Задача не найдена"
}
```

### Метод не поддерживается

Статус: `405 Method Not Allowed`

```json
{
  "error": "Метод не поддерживается"
}
```

## Архитектура

- `handlers` — HTTP хендлеры (`TasksCollection`, `TaskItem`)
- `models` — структуры данных (`Task`, `ErrorResponse`)
- `storage` — интерфейс `Storage` и in‑memory реализация `Cache` с потокобезопасностью

## Требования

- Go 1.21+
- PowerShell 7+ (для тестов на Windows) или curl (Linux/macOS)
