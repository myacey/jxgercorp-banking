# Настройка HTTPS на сервере

## 1. Сертификаты
- Используем Let's Encrypt
- Папка на проде: `/opt/certbot/jxgercorp`

## 2. Конфиг nginx
- HTTP -> HTTPS редирект
- Исключение: `/.well-known/acme-challenge/` для проверки Certbot

Пример:
```nginx
server {
    listen 80;
    server_name jxgercorp.ru www.jxgercorp.ru;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}
```

## 3. Автоматическое обновление
- Скрипт `/opt/certbot/renew.sh`:
    ```sh
    #!/bin/bash
    /usr/bin/docker run --rm \
    -v /opt/certbot/jxgercorp:/etc/letsencrypt \
    -v /opt/certbot/jxgercorp/www:/var/www/certbot \
    certbot/certbot renew --webroot -w /var/www/certbot

    if [ $? -eq 0 ]; then
    /usr/bin/docker exec frontend nginx -s reload
    fi
    ```
- Добавление задачи в `crontab`:
    ```sh
    0 0 * * * /opt/certbot/renew.sh >> /opt/certbot/renew.log 2>&1
    ```