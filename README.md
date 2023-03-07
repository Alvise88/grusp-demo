# MailUp - Demo presentazione Grusp 10 Marzo

## Run pipeline in notebook

```shell
export export HELLO_REPLICAS=3

mage demo:publish
mage demo:deploy
```

## Run pipeline from notebook

```shell
cat .argo-ci.yaml | argoci run grusp-demo - --branch="main"
```

## Smoke Test

```shell
curl -k https://hello.grusp.io/health
```