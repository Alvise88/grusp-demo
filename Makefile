SHELL := bash# we want bash behaviour in all shell invocations

# https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences
BOLD := \033[1m
NORMAL := \033[0m
GREEN := \033[1;32m

KIND_CLUSTER_NAME ?= zero
KIND_CLUSTEER_IP := $(shell ./fetch_ip.sh)

define setup_env
	./generate_dotenv.sh
endef

# # 

.PHONY: kind
bootstrap: # Create Kind cluster
	@echo "TODO: Maybe convert to a Dagger package: https://github.com/kubernetes-sigs/kind/issues/2833"
	(  kind get clusters | grep -q $(KIND_CLUSTER_NAME) ) \
	|| yq e '.networking.apiServerAddress = "$(KIND_CLUSTEER_IP)"' zero.yaml | kind create cluster --name $(KIND_CLUSTER_NAME) --config -

secrets: bootstrap
	kubectl -n cicd create secret generic kubecred --from-file=.kubeconfig=/home/vscode/.kube/config --dry-run=client -o yaml | kubectl apply -f -
	kubectl -n default create secret generic regcred --from-file=.dockerconfigjson=/home/vscode/.docker/config.json --type=kubernetes.io/dockerconfigjson --dry-run=client -o yaml | kubectl apply -f -

cluster: secrets
	kubectl kustomize --enable-alpha-plugins ./zero | kubectl apply -f -

connection: cluster
	$(call setup_env)

reset:
	kind delete cluster --name="zero"

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run -v --timeout 5m