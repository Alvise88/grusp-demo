package service

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"grusp.io/infra/imports/k8s"
)

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

func NewWebService(scope constructs.Construct, id *string, props *WebServiceProps) (constructs.Construct, error) {
	construct := constructs.NewConstruct(scope, id)

	replicas := props.Replicas
	if replicas == nil {
		replicas = jsii.Number(1)
	}

	port := props.Port
	if port == nil {
		port = jsii.Number(80)
	}

	containerPort := props.ContainerPort
	if containerPort == nil {
		containerPort = jsii.Number(8080)
	}

	label := map[string]*string{
		"app": constructs.Node_Of(construct).Id(),
	}

	svc := k8s.NewKubeService(construct, jsii.String("service"), &k8s.KubeServiceProps{
		Spec: &k8s.ServiceSpec{
			Type: jsii.String("LoadBalancer"),
			Ports: &[]*k8s.ServicePort{{
				Port:       port,
				TargetPort: k8s.IntOrString_FromNumber(containerPort),
			}},
			Selector: &label,
		},
	})

	if props.InternetFacing {
		k8s.NewKubeIngress(construct, jsii.String("ingress"), &k8s.KubeIngressProps{
			Metadata: &k8s.ObjectMeta{
				// Name: id,
				Annotations: &map[string]*string{
					"ingress.kubernetes.io/rewrite-target":        jsii.String("/$2"),
					"nginx.ingress.kubernetes.io/ssl-passthrough": jsii.String("false"),
					"kubernetes.io/ingress.class":                 jsii.String("nginx"),
				},
			},
			Spec: &k8s.IngressSpec{
				Rules: &[]*k8s.IngressRule{
					{
						Host: jsii.String(fmt.Sprintf("%s.grusp.io", *id)),
						Http: &k8s.HttpIngressRuleValue{
							Paths: &[]*k8s.HttpIngressPath{
								{
									Path:     jsii.String("/"),
									PathType: jsii.String("Prefix"),
									Backend: &k8s.IngressBackend{
										Service: &k8s.IngressServiceBackend{
											Name: svc.Name(),
											Port: &k8s.ServiceBackendPort{
												Number: port,
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
	}

	k8s.NewKubeDeployment(construct, jsii.String("deployment"), &k8s.KubeDeploymentProps{
		Spec: &k8s.DeploymentSpec{
			Replicas: replicas,
			Selector: &k8s.LabelSelector{MatchLabels: &label},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{Labels: &label},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{{
						Name:  jsii.String("web"),
						Image: props.Image,
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
								Value: props.Image,
							},
						},
						Ports: &[]*k8s.ContainerPort{{ContainerPort: containerPort}},
						LivenessProbe: &k8s.Probe{
							HttpGet: &k8s.HttpGetAction{
								Path: &props.HealthCheckPath,
								Port: k8s.IntOrString_FromNumber(containerPort),
							},
						},
						ImagePullPolicy: jsii.String("Always"),
						Resources: &k8s.ResourceRequirements{
							Limits: &map[string]k8s.Quantity{
								"memory": k8s.Quantity_FromString(jsii.String(fmt.Sprintf("%.0fMi", *props.MemoryLimit))),
								"cpu":    k8s.Quantity_FromNumber(props.CPULimit),
							},
							Requests: &map[string]k8s.Quantity{
								"memory": k8s.Quantity_FromString(jsii.String(fmt.Sprintf("%.0fMi", *props.MemoryLimit))),
								"cpu":    k8s.Quantity_FromNumber(props.CPULimit),
							},
						},
					}},
					ImagePullSecrets: &[]*k8s.LocalObjectReference{
						{
							Name: jsii.String("regcred"),
						},
					},
				},
			},
		},
	})

	return construct, nil
}
