#!/usr/bin/env bash

apk update; apk add -U --no-cache curl wget

mkdir -p /custom-tools

echo "installing mage"

mkdir -p /tmp
cd /tmp

export MAGE_VERSION="1.14.0"
curl -sLO https://github.com/magefile/mage/releases/download/v${MAGE_VERSION}/mage_${MAGE_VERSION}_Linux-64bit.tar.gz
tar -xvf mage_${MAGE_VERSION}_Linux-64bit.tar.gz
chmod +x mage
mv ./mage /custom-tools/mage
