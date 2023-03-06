#!/usr/bin/env bash

cd "$(dirname "$0")"

echo "Installing Starship"

# curl -sS https://starship.rs/install.sh | sh - -y
sh -c "$(curl -fsSL https://starship.rs/install.sh)" - -y
grep -qF "/usr/local/bin/starship" ~/.bashrc || echo "$(starship init bash)" >> ~/.bashrc

echo "Installing KSops"
export KSOPS_VERSION="3.1.1"
export XDG_CONFIG_HOME="${HOME}/.config"
grep -qF "XDG_CONFIG_HOME" ~/.bashrc || echo "export XDG_CONFIG_HOME=${HOME}/.config" >> ~/.bashrc
curl -s "https://raw.githubusercontent.com/viaduct-ai/kustomize-sops/v${KSOPS_VERSION}/scripts/install-ksops-archive.sh" | bash

grep -qF "kubectl" ~/.bashrc || echo "source <(kubectl completion bash)" >> ~/.bashrc # add autocomplete permanently to your bash shell.

# Mozilla SOPS editor configuration to use vscode
grep -qF "code --wait" ~/.bashrc ||  echo 'export EDITOR="code --wait"' >> ~/.bashrc

[ ! -f ../.env ] || source ../.env

npm install -g cdk8s-cli

# add dns entry

grep -qF "argowf.grusp.io" /etc/hosts  || echo "172.24.0.2  argowf.grusp.io" |  sudo tee -a /etc/hosts
grep -qF "hello.grusp.io" /etc/hosts  || echo "172.24.0.2   hello.grusp.io" |  sudo tee -a /etc/hosts


rm /home/vscode/.docker/config.json
../login.sh