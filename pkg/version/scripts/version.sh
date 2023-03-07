#!/usr/bin/env bash

set -e

base=${BASE}

sha=$(git rev-parse HEAD)
counter=$(git rev-list --count --no-merges HEAD)

if [[ $(git diff --stat) != '' ]]; then
  echo -n "${base}.${counter}-dirty.${sha:0:6}"
else
  echo -n "${base}.${counter}.${sha:0:6}"
fi

