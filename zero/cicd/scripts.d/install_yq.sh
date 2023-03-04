#!/usr/bin/env bash

apk update; apk add -U --no-cache curl wget

mkdir -p /custom-tools

# Install yq
echo "Installing YQ"

curl -sLo /custom-tools/yq https://github.com/mikefarah/yq/releases/download/v4.29.1/yq_linux_amd64
chmod +x /custom-tools/yq
