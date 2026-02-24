#!/bin/sh

if [ -n "$API_BASE_URL" ]; then
    echo "window.config = { apiBaseUrl: \"$API_BASE_URL\" };" > /usr/share/nginx/html/config.js
fi

exec "$@"
