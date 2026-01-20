#!/bin/bash
set -e

echo "🚀 Starting deployment..."

# Загружаем переменные
if [ -f .env.private ]; then
    echo "🔒 Loading local .env.private"
    source .env.private
else
    echo "ℹ️  .env.private not found, using environment variables"
fi

: "${SERVER_HOST:?SERVER_HOST not set}"
: "${SERVER_USER:?SERVER_USER not set}"
: "${DEPLOY_DIR:?DEPLOY_DIR not set}"
: "${DOCKERHUB_USERNAME:?DOCKERHUB_USERNAME not set}"
: "${DOCKERHUB_TOKEN:?DOCKERHUB_TOKEN not set}"

# Вход в docker
echo "🔐 Logging into DockerHub..."
echo "$DOCKERHUB_TOKEN" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin

# Сборка сервисов за исключение frontend (frontend с HMS билдится докер-флагом --profiles dev)
echo "📦 Building services..."
docker compose -f docker-compose.yml build

# Сборка фронта
echo "📦 Building frontend..."
docker compose -f docker-compose.prod.yml build frontend

# Пушим образы в docker hub
echo "⬆️ Pushing images..."
docker compose -f docker-compose.yml -f docker-compose.prod.yml push

# Отправляем конфиги на сервер
echo "📦 Packaging configs..."
tar czf deploy.tar.gz \
    docker-compose.yml \
    docker-compose.prod.yml \
    .env* \
    services/monitoring/otel-collector-config.yaml \
    services/monitoring/prometheus \
    services/monitoring/grafana

echo "📂 Sending archive to server..."
scp -o StrictHostKeyChecking=no deploy.tar.gz "${SERVER_USER}@${SERVER_HOST}:${DEPLOY_DIR}/"
rm deploy.tar.gz

# Выполняем деплой по ssh
echo "🚀 Running remote deployment..."
ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_HOST} << EOF
    cd ${DEPLOY_DIR}
    tar xzf deploy.tar.gz
    rm deploy.tar.gz
    sudo docker compose -f docker-compose.yml -f docker-compose.prod.yml pull
    sudo docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --remove-orphans --force-recreate
EOF

echo "✅ Deployment completed successfully!"
