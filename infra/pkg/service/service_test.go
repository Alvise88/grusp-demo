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
				Image:           jsii.String("ghost"),
				InternetFacing:  false,
				Port:            jsii.Number(2368),
				ContainerPort:   jsii.Number(2368),
				HealthCheckPath: "/health",
				MemoryLimit:     jsii.Number(64),
				CPULimit:        jsii.Number(0.1),
				Replicas:        jsii.Number(1),
			},
		},
		{
			props: WebServiceProps{
				Image:           jsii.String("ghost"),
				InternetFacing:  true,
				Port:            jsii.Number(2368),
				ContainerPort:   jsii.Number(2368),
				HealthCheckPath: "/health",
				MemoryLimit:     jsii.Number(64),
				CPULimit:        jsii.Number(0.1),
				Replicas:        jsii.Number(1),
			},
		},
	}

	for _, tc := range cases {
		chart := cdk8s.Testing_Chart()

		_, err := NewWebService(chart, jsii.String("test"), &tc.props)

		if err != nil {
			t.Error(err)
			return
		}

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
