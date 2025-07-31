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
- `make up` — быстро поднять dev-среду: запускает PostgreSQL в Docker, применяет миграции, запускает приложение с тегом `dev`  
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

## Тесты и покрытие

- Полностью покрыты тестами транспортный и сервисный слои  
- Все зависимости мокируются с помощью `mockery` для изоляции компонентов  
- Отчёт по покрытию тестов можно получить командой `make go-tests-coverage` — откроется в браузере  

---
