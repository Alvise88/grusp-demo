# MailUp - Demo presentazione Grusp 10 Marzo

Utilizzo di Dagger per una ci/cd pipeline legata a K8s. 

## Esempio

"Architettura" dell'applicazione che verrà installata nel cluster K8s.

![APP](./docs/img/hello-grusp.excalidraw.png)

Possibile astrazione da fornire al team di sviluppo:

```go
type WebServiceProps struct {
	constructs.Construct
	Image           *string
	InternetFacing  bool
	Port            *float64
	ContainerPort   *float64
	HealthCheckPath string
	MemoryLimit     *float64
	CPULimit        *float64
	Replicas        *float64
}
```

Attuale configurazione del servizio:

```go
service.NewWebService(chart, jsii.String("hello"), &service.WebServiceProps{
	Image:           jsii.String(fmt.Sprintf("alvisevitturi/hello-grusp:%s", version)),
	InternetFacing:  true,
	Port:            jsii.Number(8080),
	ContainerPort:   jsii.Number(8080),
	HealthCheckPath: "/health",
	MemoryLimit:     jsii.Number(64),
	CPULimit:        jsii.Number(0.1),
	Replicas:        jsii.Number(replicas),
})
```

## Run pipeline in notebook

```shell
export export HELLO_REPLICAS=3

mage demo:test
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