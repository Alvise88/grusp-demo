#!/usr/bin/env bash

apk update; apk add -U --no-cache curl wget

mkdir -p /custom-tools

echo "installing jq"

# Install jq
curl -sLo /custom-tools/jq  https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
chmod +x /custom-tools/jq
