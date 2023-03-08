package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"grusp.io/infra/pkg/service"
)

type GruspChartProps struct {
	cdk8s.ChartProps
}

func NewGruspChart(scope constructs.Construct, id string, props *GruspChartProps) (cdk8s.Chart, error) {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	version, ok := os.LookupEnv("HELLO_VERSION")

	if !ok {
		panic(errors.New("hello version must be provided"))
	}

	envReplicas, ok := os.LookupEnv("HELLO_REPLICAS")

	if !ok || envReplicas == "" {
		envReplicas = "1"
	}

	replicas, err := strconv.ParseFloat(envReplicas, 64)

	if err != nil {
		return nil, err
	}

	_, err = service.NewWebService(chart, jsii.String("hello"), &service.WebServiceProps{
		Image:           jsii.String(fmt.Sprintf("alvisevitturi/hello-grusp:%s", version)),
		InternetFacing:  true,
		Port:            jsii.Number(8080),
		ContainerPort:   jsii.Number(8080),
		HealthCheckPath: "/health",
		MemoryLimit:     jsii.Number(64),
		CPULimit:        jsii.Number(0.1),
		Replicas:        jsii.Number(replicas),
	})

	return chart, err
}

func main() {
	app := cdk8s.NewApp(nil)

	_, err := NewGruspChart(app, "hello", nil)

	if err != nil {
		panic(err)
	}

	app.Synth()
}
