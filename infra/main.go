package main

import (
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

	envReplicas, ok := os.LookupEnv("HELLO_REPLICAS")

	if !ok {
		envReplicas = "1"
	}

	replicas, err := strconv.ParseFloat(envReplicas, 64)

	if err != nil {
		return nil, err
	}

	service.NewWebService(chart, jsii.String("hello"), &service.WebServiceProps{
		Image:          jsii.String("alvisevitturi/hello-grusp:latest"),
		Port:           jsii.Number(8080),
		ContainerPort:  jsii.Number(8080),
		InternetFacing: true,

		Replicas: jsii.Number(replicas),
	})

	return chart, nil
}

func main() {
	app := cdk8s.NewApp(nil)

	_, err := NewGruspChart(app, "hello", nil)

	if err != nil {
		panic(err)
	}

	app.Synth()
}
