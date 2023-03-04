#!/usr/bin/env bash

apk update; apk add -U --no-cache curl wget

mkdir -p /custom-tools

# Install dagger
# Could be removed by using an official muli-arch dagger.io image
arch=$(uname -m)
case $arch in
  x86_64) arch="amd64" ;;
  aarch64) arch="arm64" ;;
esac

[[ -z "${DAGGER_VERSION}" ]] && dagger_version='0.2.36' || dagger_version=${DAGGER_VERSION}

wget -c https://github.com/dagger/dagger/releases/download/v${dagger_version}/dagger_v${dagger_version}_linux_${arch}.tar.gz  -O - |  \
tar zxf - -C /custom-tools/
chmod +x /custom-tools/dagger
