#!/bin/sh

if [ -n "$API_BASE_URL" ]; then
    # Use jq to safely encode the API_BASE_URL as a JSON string
    ENCODED_API_BASE_URL=$(jq -n --arg url "$API_BASE_URL" '$url')
    echo "window.config = { apiBaseUrl: $ENCODED_API_BASE_URL };" > /usr/share/nginx/html/config.js
fi

exec "$@"
