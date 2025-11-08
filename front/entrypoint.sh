#!/bin/sh
set -e

# Генерируем конфиг
if [ "$NODE_ENV" = "development" ]; then
    echo "🚀 Development mode: starting Vite dev server..."

    export CHOKIDAR_USEPOLLING=true
    export WATCHPACK_POLLING=true

    exec yarn serve --host 0.0.0.0 --port 80
elif [ "$NODE_ENV" = "prod" ]; then
    echo "🏗️ Production mode: serving built files..."

    CONFIG_FILE=/usr/share/nginx/html/config.js
    echo "window.__APP_CONFIG__ = {" > $CONFIG_FILE
    echo "  VITE_API_BASE_URL: '${VITE_API_BASE_URL:-/api}'" >> $CONFIG_FILE
    echo "};" >> $CONFIG_FILE

    echo "Generated runtime config:"
    cat $CONFIG_FILE

    exec nginx -g 'daemon off;'
fi
