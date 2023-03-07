#!/usr/bin/env bash

# shellcheck disable=SC2016
set -e

# apk update; apk add -U --no-cache curl wget

mkdir -p /custom-tools

VERSION="1.20.1"

# rm -rf /custom-tools/go && tar -C /custom-tools -xzf go${VERSION}.linux-amd64.tar.gz

OS="$(uname -s)"
ARCH="$(uname -m)"

case $OS in
    "Linux")
        case $ARCH in
        "x86_64")
            ARCH=amd64
            ;;
        "aarch64")
            ARCH=arm64
            ;;
        "armv6" | "armv7l")
            ARCH=armv6l
            ;;
        "armv8")
            ARCH=arm64
            ;;
        .*386.*)
            ARCH=386
            ;;
        esac
        PLATFORM="linux-$ARCH"
    ;;
    "Darwin")
          case $ARCH in
          "x86_64")
              ARCH=amd64
              ;;
          "arm64")
              ARCH=arm64
              ;;
          esac
        PLATFORM="darwin-$ARCH"
    ;;
esac

PACKAGE_NAME="go$VERSION.$PLATFORM.tar.gz"

rm -rf /custom-tools/go && wget -c https://storage.googleapis.com/golang/$PACKAGE_NAME -O - |  \
tar zxf - -C /custom-tools/