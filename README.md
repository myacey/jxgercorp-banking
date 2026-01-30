# JXGERcorp Banking
**JXGERcorp Banking** — это транзакционный приватный банкинг для бизнеса и физических лиц, реализованный на Go. Проект представляет собой микросервисную архитектуру с использованием современных технологий: REST API, gRPC, Kafka, Jaeger для трассировки, PostgreSQL для хранения данных, Redis для управления сессиями, и Vue для фронтенда.


[ДОБАВИТЬ k8s]: #

<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=go,postgres,redis,vue,docker,kafka,gmail,github,prometheus,grafana,jaeger" />
  </a>
</p>

![app-architecture](https://github.com/user-attachments/assets/fa19ff84-b898-4274-a752-ed31eb87738c)



## Описание
JXGERcorp Banking предоставляет пользователям следующие возможности:
1. **Регистрация и подтверждение учетной записи**  
    Пользователь создаётся с статусом `pending=true` (ожидает подтверждения). Генерируется код подтверждения, который через `Kafka` отправляется на email.
2. **Авторизация и вход в систему**  
    Пользователь получает токен (`JWT`), который используется для доступа к защищенным ресурсам через `Bearer Token` и/или cookie `authToken`.
3. **Управление аккаунтами (картами)**  
    У пользователя может быть несколько аккаунтов в разных валютах (RUB, USD, EUR). Каждый аккаунт создается с балансом 1000 единиц соответствующей валюты. Пользователь может создавать новые аккаунты и удалять существующие.
4. **Создание транзакций**  
    Перевод средств между аккаунтами выполняется через `REST API`. Все операции атомарные, а баланс аккаунта не может быть отрицательным.
5. **Просмотр истории операций и аккаунтов**  
    Возможность отслеживать состояние всех своих аккаунтов. История транзакций для аккаунтов поддерживает постраничный просмотр.

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

- **GET** `/api/v1/user/confirm?username=string&code=string`:
    - Подтверждает учетную запись по коду из email

### API Transfer/Account
- **POST** `/api/v1/transfer/account`
    - Создание нового аккаунта (карты) пользователя.
    - Тело запроса:
    ```json
    {
        "currency": "RUB"
    }
    ```
    - Пример ответа:
    ```json
    {
        "id": "536060ea-df8d-4de4-bc92-1dafb0247f4f",
        "owner_username": "myacey",
        "balance": 1000,
        "currency": "EUR",
        "created_at": "2026-01-30T13:15:22.359059Z"
    }
    ```

- **DELETE** `/api/v1/transfer/account?account_id=string`
    - Удаление аккаунта пользователя.
    - Пример ответа:
    ```json
    {
        "message": "success"
    }
    ```

- **GET** `/api/v1/transfer/accounts?username=string&currency=string`
    - Получение списка аккаунтов пользователя, можно фильтровать по валюте.
    - Пример ответа:
    ```json
    [
        {
            "id": "536060ea-df8d-4de4-bc92-1dafb0247f4f",
            "owner_username": "myacey",
            "balance": 1000,
            "currency": "EUR",
            "created_at": "2026-01-30T13:15:22.359059Z"
        },
        {
            "id": "966a13f6-3634-4813-b036-5daab33b1f3a",
            "owner_username": "myacey",
            "balance": 1000,
            "currency": "RUB",
            "created_at": "2026-01-30T13:18:04.037894Z"
        }
    ]
    ```

### API Transfer / Transactions

- **POST** `/api/v1/transfer`
    - Создание перевода между аккаунтами.
    - Тело запроса:
    ```json
    {
        "from_account_id":"2d4fefd4-98b5-42e8-af22-efd2bdc600da",
        "from_account_username":"ritzwi",
        "to_account_id":"4b22d369-f505-4646-af92-768a9854e6dd",
        "to_account_username":"myacey",
        "amount":9,
        "currency":"EUR"
    }
    ```
    - Пример ответа:
    ```json
    {
        "id": "54c35fc0-e4eb-423b-89ac-d84a9212be0a",
        "from_account_id": "2d4fefd4-98b5-42e8-af22-efd2bdc600da",
        "from_account_username": "ritzwi",
        "to_account_id": "4b22d369-f505-4646-af92-768a9854e6dd",
        "to_account_username": "myacey",
        "amount": 9,
        "currency": "EUR",
        "created_at": "2026-01-30T13:26:06.543936Z"
    }
    ```

- **GET** `/api/v1/transfer?current_account_id=string`
    - Получение истории транзакций для конкретного аккаунта.
    - Пример ответа:
    ```json
    [
        {
            "id": "54c35fc0-e4eb-423b-89ac-d84a9212be0a",
            "from_account_id": "2d4fefd4-98b5-42e8-af22-efd2bdc600da",
            "from_account_username": "ritzwi",
            "to_account_id": "4b22d369-f505-4646-af92-768a9854e6dd",
            "to_account_username": "myacey",
            "amount": 9,
            "currency": "EUR",
            "created_at": "2026-01-30T13:26:06.543936Z"
        },
        {
            "id": "40f366dc-9552-4099-af88-d73e7b17c344",
            "from_account_id": "4b22d369-f505-4646-af92-768a9854e6dd",
            "from_account_username": "myacey",
            "to_account_id": "2d4fefd4-98b5-42e8-af22-efd2bdc600da",
            "to_account_username": "ritzwi",
            "amount": 100,
            "currency": "EUR",
            "created_at": "2026-01-30T13:25:39.741689Z"
        }
    ]
    ```

- **GET** `/api/v1/transfer/currencies`
    - Получение списка поддерживаемых валют.
    - Пример ответа:
    ```json
    [
        {
            "code": "RUB",
            "symbol": "₽",
            "precision": 2
        },
        {
            "code": "USD",
            "symbol": "$",
            "precision": 2
        },
        {
            "code": "EUR",
            "symbol": "€",
            "precision": 2
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
- **OTEL Collector** — сбор метрик с сервисов в единой точке, дальнейший экспорт
- **Grafana и Prometheus** — мониторинг состояния сервисов и бизнес-метрик
- **SMTP** — отправка почты

### Frontend:
- **Vue** — SPA для работы пользователя



## Архитектура и структура проекта
Проект организован в виде монорепозитория, разделенного на две основные части:
- **front/**  — фронтенд приложения на Vue.
- **services/** — микросервисы на Go:
    - api-gateway — точка входа в систему, проксирует запросы, обрабатывает авторизацию и трассировку.
    - user — управление пользователями (регистрация, авторизация, получение баланса и подтверждение email).
    - transfer — обработка переводов и истории транзакций.
    - token — управление токенами и аутентификацией (gRPC).
    - notification — сервис отправки email уведомлений, взаимодействующий через Kafka.
    - libs — общие библиотеки и протоколы.
    - monitoring — конфигурация observability

Каждый микросервис состоит из слоев:
- **Controller** — обработка HTTP-запросов, валидация, логирование
- **Service** — бизнес логика: логика пользователей, транзакций, запросы в другие микросервисы (`gRPC`, `Kafka`)
- **Repository** — работа с базой данных: управления данными пользователей, операциями

## Установка и запуск

1. **Клонируйте репозиторий**:
```sh
git clone https://github.com/myacey/jxgercorp-banking.git
cd jxgercorp-banking
```

2. **Создайте `.env.private` файл**:
```sh
# MAIL
SMTP_MAIL_ADRESS='your_mail_adress'
SMTP_PASSWORD='your_password'
```
> [!TIP]
> Необходимо для проверки почты при регистрации

2. **Запустите сервис локально**:
```sh
make up
```

3. **Задеплойте сервис в облако**:

    3.1. **Добавьте в `.env.private`**:
    ```sh
    SERVER_USER=your_server_username
    SERVER_HOST=your_server_host
    DEPLOY_DIR=/your/deploy/dir
    DOCKERHUB_USERNAME=your_dockerhub_username
    DOCKERHUB_TOKEN=your_dockerhub_token
    ```

    3.2. **Запустите скрипт деплоя**
    ```sh
    bash ./scripts/deploy-prod.sh
    ```
> [!CAUTION]
> Это запустит деплой на указанный вами прод. Убедитесь, что на нем уже есть .env.private файл.
