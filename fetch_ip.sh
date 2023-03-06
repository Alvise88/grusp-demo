#!/usr/bin/env sh

set -e

ifconfig eth0 | awk '/inet / {print $2; }' | cut -d ' ' -f 2 |  tr -d '\n'