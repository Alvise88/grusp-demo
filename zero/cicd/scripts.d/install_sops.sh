#!/usr/bin/env bash

apk update; apk add -U --no-cache curl wget

mkdir -p /custom-tools

echo "installing sops"

# Install Sops
[[ -z "${SOPS_VERSION}" ]] && sops_version='3.7.3' || sops_version=${SOPS_VERSION}

curl -sLo /custom-tools/sops https://github.com/mozilla/sops/releases/download/v${sops_version}/sops-v${sops_version}.linux.amd64
chmod +x /custom-tools/sops
