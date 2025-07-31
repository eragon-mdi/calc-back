# Calc-back Microservice

Микросервис для работы с [вычислениями](https://github.com/EvgenyMentor/CalculatorAppFrontendPantela).

## API ручки

| Метод  | Путь                   | Описание                               |
|--------|------------------------|---------------------------------------|
| GET    | `/calculations`         | Получить последние 10 вычислений      |
| GET    | `/calculations/{id}`    | Получить вычисление по UUID            |
| POST   | `/calculations`         | Создать новое вычисление по выражению |
| DELETE | `/calculations/{id}`    | Удалить вычисление по UUID             |
| PATCH  | `/calculations/{id}`    | Обновить вычисление по UUID            |

Полное описание доступно в Swagger-документации.

Swagger можно развернуть локально, запустив Go-приложение с тегом `dev` (как указано в Makefile), и открыть в браузере по адресу: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

---

## Быстрый старт

В Makefile предусмотрены команды для удобной работы:

- `make start-quiet` — продовый запуск  
- `make dev-up` — быстро поднять dev-среду: запускает PostgreSQL в Docker, применяет миграции, запускает приложение с тегом `dev`  
- `make migrate-new name=имя_миграции` — создать новую миграцию с указанным именем  

Дополнительно доступны:

- `make migrate-up`  
- `make lint`  
- `make swag`  
- `make gen-mocks`  
- `make gen-base-tests-transport`  
- `make gen-base-tests-service`  
- `make go-tests-coverage`  

---

## Конфигурация

Пример основных переменных окружения (env):

```env
STORAGE_HOST=db
STORAGE_PORT=5432
STORAGE_USER=postgres
STORAGE_PASS=secret
STORAGE_NAME=mydb
STORAGE_SSLM=disable

SERVER_ADDR=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_READ_HEADER_TIMEOUT=5s
SERVER_IDLE_TIMEOUT=60s

LOGGER_LEVEL=info
LOGGER_ENCODING=json
LOGGER_OUTPUT=stdout
LOGGER_MESSAGE_KEY=message

MIDDLEWARE_AUTH_TOKEN=dsakdjaskjkj
```
>Важно:
>- Для локальной разработки без контейнеров в STORAGE_HOST ставится `localhost`.
>- Для запуска в Docker Compose — STORAGE_HOST должен быть равен имени сервиса базы в оркестраторе: `db`.

## Тесты и покрытие

- Полностью покрыты тестами транспортный и сервисный слои  
- Все зависимости мокируются с помощью `mockery` для изоляции компонентов  
- Отчёт по покрытию тестов можно получить командой `make go-tests-coverage` — откроется в браузере  

---
