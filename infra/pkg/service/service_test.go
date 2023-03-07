package service

import (
	"testing"

	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewWebService(t *testing.T) {
	type testCase struct {
		props WebServiceProps
	}

	cases := []testCase{
		{
			props: WebServiceProps{
				Image:          jsii.String("ghost"),
				Replicas:       jsii.Number(1),
				ContainerPort:  jsii.Number(2368),
				InternetFacing: false,
			},
		},
		{
			props: WebServiceProps{
				Image:          jsii.String("ghost"),
				Replicas:       jsii.Number(1),
				ContainerPort:  jsii.Number(2368),
				InternetFacing: true,
			},
		},
	}

	for _, tc := range cases {
		chart := cdk8s.Testing_Chart()

		NewWebService(chart, jsii.String("test"), &tc.props)

		manifests := cdk8s.Testing_Synth(chart)

		ingresses := ([]map[string]interface{}{})
		services := ([]map[string]interface{}{})
		deployments := ([]map[string]interface{}{})

		for _, mnf := range *manifests {
			manifest := mnf.(map[string]interface{})

			if manifest["kind"] == "Ingress" {
				ingresses = append(ingresses, manifest)
			}

			if manifest["kind"] == "Service" {
				services = append(services, manifest)
			}

			if manifest["kind"] == "Deployment" {
				deployments = append(deployments, manifest)
			}
		}

		if tc.props.InternetFacing {
			assert.Equal(t, len(ingresses), 1)
		} else {
			assert.Equal(t, len(ingresses), 0)
		}

		assert.Equal(t, len(deployments), 1)
		assert.Equal(t, len(services), 1)
	}
}
