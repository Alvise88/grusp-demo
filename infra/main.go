package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"grusp.io/infra/imports/k8s"
)

type GruspChartProps struct {
	cdk8s.ChartProps
}

func NewGruspChart(scope constructs.Construct, id string, props *GruspChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	label := map[string]*string{"app": jsii.String("hello-grusp")}

	k8s.NewKubeService(chart, jsii.String("hello-grusp"), &k8s.KubeServiceProps{
		Spec: &k8s.ServiceSpec{
			Type: jsii.String("LoadBalancer"),
			Ports: &[]*k8s.ServicePort{{
				Port:       jsii.Number(18080),
				TargetPort: k8s.IntOrString_FromNumber(jsii.Number(8080)),
			}},
			Selector: &label,
		},
	})

	k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(2),
			Selector: &k8s.LabelSelector{
				MatchLabels: &label,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &label,
				},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{{
						Name:  jsii.String("hello-grusp"),
						Image: jsii.String("alvisevitturi/hello-grusp:latest"),
						Ports: &[]*k8s.ContainerPort{{ContainerPort: jsii.Number(8080)}},
					}},
				},
			},
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewGruspChart(app, "hello", nil)
	app.Synth()
}
