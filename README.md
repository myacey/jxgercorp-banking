# JXGERcorp Banking
**JXGERcorp Banking** — это транзакционный приватный банкинг для бизнеса и физических лиц, реализованный на Go. Проект представляет собой микросервисную архитектуру с использованием современных технологий: REST API, gRPC, Kafka, Jaeger для трассировки, PostgreSQL для хранения данных, Redis для управления сессиями, и Vue для фронтенда.


[ДОБАВИТЬ KAFKA, GRAFANA, PROMETHEUS, 3S]: #

<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=go,postgres,redis,vue,docker,kafka,gmail,github" />
  </a>
</p>

## Описание
JXGERcorp Banking предоставляет пользователям следующие возможности:
1. **Регистрация и подтверждение учетной записи**  
   При регистрации пользователь создаётся со статусом `pending=true` (ожидает подтверждения). После этого генерируется код подтверждения, и через Kafka-уведомления отправляется письмо с ссылкой для подтверждения почты.
2. **Авторизация и вход в систему**  
   Пользователь получает токен (PASETO), который используется для доступа к защищенным ресурсам через Bearer Token и/или cookie `authToken`.
3. **Создание транзакций**  
   Сервис позволяет переводить средства между пользователями через REST API, а транзакции обрабатываются атомарно. Каждый пользователь при регистрации получает начальный баланс (*1000\$*), а баланс никогда не может быть отрицательным.
4. **Просмотр истории операций**  
   История транзакций доступна с возможностью постраничного просмотра (через параметры `offset` и `limit`).

Каждый пользователь при регистрации получает 1000\$. Баланс не может быть отрицательным, а все операции выполняются атомарно.

## Основные функции API

[ПЕРЕНЕСТИ В SWAGGER]: #
### API User
- **POST** `/api/v1/user/register`:
    - Регистрация пользователя в сервисе.
    - Тело запроса:
        ```json
        {
            "username": "string",
            "email": "string",
            "password": "string"
        }
        ```
    - Пример ответа:
        ```json
        {
            "message": "success"
        }
        ```
    - После регистрации пользователь получает статус pending, и уведомление для подтверждения отправляется через Kafka.

- **POST** `/api/v1/user/login`:
    - Получения токена для доступа к системе.
    - Тело запроса:
        ```json
        {
            "username": "string",
            "password": "string"
        }
        ```
    - Пример ответа:
        ```json
        {
            "message": "success"
        }
        ```
        > Также устанавливается **Cookie** `authToken`

- **GET** `/api/v1/user/balance`:
    - Получение баланса пользователя.
    - Заголовок запроса:
        ```sh
        Bearer Token: your-auth-token
        ```
    - Пример ответа:
        ```json
        {
            "balance": 1000
        }
        ```

- **POST** `/api/v1/user/confirm`:
    - Подтверждает учетную запись по коду из email
    - Тело запроса:
        ```json
        {
            "username": "string",
            "code": "string"
        }
        ```
    - Пример ответа:
        ```json
        {
            "message": "success"
        }
        ```

### API Transaction
- **POST** `api/v1/transaction/create`:
    - Создает транзакцию (перевод денег между пользователями).
    - Тело запроса:
        ```json
        {
            "to_user": "string",
            "amount": 100
        }
        ```
        > Запрос защищен: требуется `Bearer Token`.
    - Пример ответа:
        ```json
        {
            "id": 1,
            "from_user": "string",
            "to_user": "string",
            "amount": 100,
            "created_at": "2025-02-08T11:41:24.425355Z"
        }
        ```

- **GET** `api/v1/transaction/search`:
    - Получает историю транзакций с поддержкой постраничного просмотра.
    - Тело запроса:
        ```json
        {
            "offset": 0,
            "limit": 11
        }
        ```
        > Если записей больше 10, отображаются только первые 10, а кнопка "Next" становится активной.
    - Пример ответа:
        ```json
        [
            {
                "with_user": "userB",
                "amount": 100,
                "created_at": "2025-02-08T11:41:24.425355Z"
            },
            {
                "with_user": "userC",
                "amount": -50,
                "created_at": "2025-02-08T12:15:00.000000Z"
            }
        ]
        ```

## Технологический стек:

### Backend:
- **Go** — основный язык программирования
- **PostgreSQL** — основное хранилище данных
- **Redis** — для управления сессиями
- **Kafka** — брокер сообщений для уведомлений (например, подтверждения email)
- **gRPC** — для межсервисного взаимодействия (например, сервис токенов)
- **Docker** — для быстрой упаковки и запуска
- **Jaeger** — для распределенной трассировки запросов
- **SMTP** — отправка почты
- **Grafana и Prometheus** (совсем скоро)...

### Frontend:
- **Vue** — быстрый и простой фронтенд фреймворк



## Архитектура и структура проекта
Проект организован в виде монорепозитория, разделенного на две основные части:
- **front/**  — фронтенд приложения на Vue.
- **services/** — микросервисы на Go:
    - api-gateway — точка входа в систему, проксирует запросы, обрабатывает авторизацию и трассировку.
    - user — управление пользователями (регистрация, авторизация, получение баланса и подтверждение email).
    - transaction — обработка переводов и истории транзакций.
    - token — управление токенами и аутентификацией (gRPC).
    - notification — сервис отправки email уведомлений, взаимодействующий через Kafka.
    - db — скрипты миграций, sqlc-код и другие файлы, связанные с базой данных.
    - shared — общие библиотеки, конфигурации, ошибки, протоколы (protobuf) и утилиты.

Каждый микросервис состоит из слоев:
- **Controller** — обработка HTTP-запросов, валидация, логирование
- **Service** — бизнес логика: логика пользователей, транзакций, запросы в другие микросервисы (`gRPC`, `Kafka`)
- **Repository** — работа с базой данных: управления данными пользователей, операциями

## Установка и запуск

### Локальный запуск

1. **Клонируйте репозиторий**:
```sh
git clone https://github.com/myacey/jxgercorp-banking.git
cd jxgercorp-banking
```

2. **Создайте `.env` файл в корне**:
```sh
APP_DOMAIN=http://localhost:8080
DB_PASSWORD=secret

# POSTGRES
POSTGRES_HOST=localhost
POSTGRES_USER=root
POSTGRES_DB=jxger_bank
POSTGRES_PORT=5432

# REDIS
REDIS_USER=root
REDIS_ADDRESS=redis:6379

# MAIL
GOOGLE_MAIL_ADRESS='yourmail@example.com'
GOOGLE_APP_PASSWORD='your password here'
```

3. **Запустите интеграционные сервисы в Docker**:
```sh
docker compose up -d
```
Это запустит все инфраструктурные сервисы: PostgreSQL, Redis, Kafka, Jaeger и другие.


4. **Запустите фронтенд**
```sh
cd front
npm install
npm run serve
```

4. **Запустите микросервисы**: Из корневой директории запускайте каждый микросервис в отдельных терминалах:
```sh
go run services/api-gateway/cmd/main.go
go run services/user/cmd/main.go
go run services/transaction/cmd/main.go
go run services/token/cmd/main.go
go run services/notification/cmd/main.go
```

## Дополнительные задачи (в процессе разработки)
- [ ] Реализация **Refresh Tokens** и **API Tokens**.
- [ ] Интеграция **Prometheus** и **Grafana** для метрик и алертов.
- [ ] Развертывание в **Kubernetes** и облачных платформах (например, **Yandex Cloud** или **AWS S3**).
- [ ] Автодеплой через **GitHub Actions**.

## Заключение
**JXGERcorp Banking** — это современное решение для транзакционного банкинга, построенное на микросервисной архитектуре с использованием *Go* и *Vue*. Проект включает полноценную поддержку авторизации, распределенную трассировку с *Jaeger*, обмен сообщениями через *Kafka* и масштабируемую инфраструктуру с *Docker*.
