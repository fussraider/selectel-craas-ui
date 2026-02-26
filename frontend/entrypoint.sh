#!/bin/sh

CSP_CONNECT_SRC=""

if [ -n "$API_BASE_URL" ]; then
    # Use jq to safely encode the API_BASE_URL as a JSON string
    ENCODED_API_BASE_URL=$(jq -n --arg url "$API_BASE_URL" '$url')
    echo "window.config = { apiBaseUrl: $ENCODED_API_BASE_URL };" > /usr/share/nginx/html/config.js

    # Extract origin if API_BASE_URL is absolute
    if echo "$API_BASE_URL" | grep -q "^http"; then
        CSP_CONNECT_SRC=$(echo "$API_BASE_URL" | sed -E 's|^(https?://[^/]+).*|\1|')
    fi
fi

# Replace CSP connect source placeholder in nginx config
# using | as delimiter since it is not valid in URLs (RFC 3986)
sed -i "s|__CSP_CONNECT_SRC__|$CSP_CONNECT_SRC|g" /etc/nginx/conf.d/default.conf

exec "$@"
