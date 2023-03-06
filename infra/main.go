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

	// - appProtocol: https
	//   name: https
	//   port: 443
	//   protocol: TCP
	//   targetPort: https

	// k8s.NewKubeService(chart, jsii.String("hello-grusp"), &k8s.KubeServiceProps{
	// 	Spec: &k8s.ServiceSpec{
	// 		Type: jsii.String("NodePort"),
	// 		Ports: &[]*k8s.ServicePort{{
	// 			AppProtocol: jsii.String("http"),
	// 			Name:        jsii.String("http"),
	// 			Port:        jsii.Number(8080),
	// 			NodePort:    jsii.Number(31496),
	// 			Protocol:    jsii.String("TCP"),
	// 			TargetPort:  k8s.IntOrString_FromString(jsii.String("http")), // jsii.String("http"), // k8s.IntOrString_FromNumber(jsii.Number(8080)),
	// 		}},
	// 		Selector: &label,
	// 	},
	// })

	// k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
	// 	Spec: &k8s.DeploymentSpec{
	// 		Replicas: jsii.Number(2),
	// 		Selector: &k8s.LabelSelector{
	// 			MatchLabels: &label,
	// 		},
	// 		Template: &k8s.PodTemplateSpec{
	// 			Metadata: &k8s.ObjectMeta{
	// 				Labels: &label,
	// 			},
	// 			Spec: &k8s.PodSpec{
	// 				Containers: &[]*k8s.Container{{
	// 					Name:  jsii.String("hello-grusp"),
	// 					Image: jsii.String("alvisevitturi/hello-grusp:latest"),
	// 					Ports: &[]*k8s.ContainerPort{
	// 						{
	// 							ContainerPort: jsii.Number(8080),
	// 							// HostPort:      jsii.Number(18080),
	// 							Name:     jsii.String("http"),
	// 							Protocol: jsii.String("TCP"),
	// 						},
	// 					},
	// 				}},
	// 			},
	// 		},
	// 	},
	// })

	svc := k8s.NewKubeService(chart, jsii.String("hello-grusp"), &k8s.KubeServiceProps{
		Spec: &k8s.ServiceSpec{
			Type: jsii.String("LoadBalancer"),
			Ports: &[]*k8s.ServicePort{{
				Port:       jsii.Number(8080),
				TargetPort: k8s.IntOrString_FromNumber(jsii.Number(8080)),
			}},
			Selector: &label,
		},
	})

	k8s.NewKubeIngress(chart, jsii.String("ingress"), &k8s.KubeIngressProps{
		Metadata: &k8s.ObjectMeta{
			Name: jsii.String("hello-ingress"),
			Annotations: &map[string]*string{
				"ingress.kubernetes.io/rewrite-target":        jsii.String("/$2"),
				"nginx.ingress.kubernetes.io/ssl-passthrough": jsii.String("false"),
				"kubernetes.io/ingress.class":                 jsii.String("nginx"),
			},
		},
		Spec: &k8s.IngressSpec{
			Rules: &[]*k8s.IngressRule{
				{
					Host: jsii.String("hello.grusp.io"),
					Http: &k8s.HttpIngressRuleValue{
						Paths: &[]*k8s.HttpIngressPath{
							{
								Path:     jsii.String("/"),
								PathType: jsii.String("Prefix"),
								Backend: &k8s.IngressBackend{
									Service: &k8s.IngressServiceBackend{
										Name: svc.Name(),
										Port: &k8s.ServiceBackendPort{
											Number: jsii.Number(8080),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	k8s.NewKubeDeployment(chart, jsii.String("hello-deployment"), &k8s.KubeDeploymentProps{
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
