#!/usr/bin/env bash

set -e

echo "Install cdk8s"

K8S_VERSION="1.23.16"

mkdir -p /tmp/k8s

cd  /tmp/k8s
curl -LO https://dl.k8s.io/release/v${K8S_VERSION}/bin/linux/amd64/kubectl
install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
cd /tmp

rm -R  /tmp/k8s

npm install -g cdk8s-cli

ln -s /usr/local/share/nvm/versions/node/v18.14.2/bin/cdk8s /usr/local/bin/cdk8s