#!/usr/bin/env bash

set -eu

mkdir -p /custom-tools

export KUBECTL_VERSION="1.25.3"
curl -LO https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl
chmod +x kubectl
mv kubectl /custom-tools/kubectl
