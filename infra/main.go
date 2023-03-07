package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"grusp.io/infra/imports/k8s"
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

	imageName := fmt.Sprintf("alvisevitturi/hello-grusp:%s", version)

	service.NewWebService(chart, jsii.String("hello"), &service.WebServiceProps{
		Image:         jsii.String(imageName),
		Port:          jsii.Number(8080),
		ContainerPort: jsii.Number(8080),
		Env: &[]*k8s.EnvVar{
			{
				Name: jsii.String("KUBERNETES_NAMESPACE"),
				ValueFrom: &k8s.EnvVarSource{
					FieldRef: &k8s.ObjectFieldSelector{FieldPath: jsii.String("metadata.namespace")},
				},
			},
			{
				Name: jsii.String("KUBERNETES_NODE_NAME"),
				ValueFrom: &k8s.EnvVarSource{
					FieldRef: &k8s.ObjectFieldSelector{FieldPath: jsii.String("spec.nodeName")},
				},
			},
			{
				Name: jsii.String("KUBERNETES_POD_NAME"),
				ValueFrom: &k8s.EnvVarSource{
					FieldRef: &k8s.ObjectFieldSelector{FieldPath: jsii.String("metadata.name")},
				},
			},
			{
				Name:  jsii.String("CONTAINER_IMAGE"),
				Value: jsii.String(imageName),
			},
		},

		AlwaysPull:     true,
		InternetFacing: true,

		Replicas: jsii.Number(replicas),
		Health:   "/health",
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
