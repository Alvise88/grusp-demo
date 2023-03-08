package k8s

import (
	"fmt"
	"strings"

	"dagger.io/dagger"
	"github.com/alvise88/grusp-demo/pkg/debian"
)

type Cdk8sOpts struct {
	Image   *dagger.Container
	Version string
}

func Cdk8s(c *dagger.Client, opts Cdk8sOpts) (*dagger.Container, error) {
	base := opts.Image

	if base == nil {
		debian, err := debian.Debian(c, debian.DebianOpts{
			Version: "11.6",
			Packages: []struct {
				Name    string
				Version string
			}{
				{
					Name: "bash",
				},
				{
					Name: "apt-utils",
				},
				{
					Name: "curl",
				},
				{
					Name: "wget",
				},
				{
					Name: "openssl",
				},
				{
					Name: "ca-certificates",
				},
			},
		})

		if err != nil {
			return nil, err
		}

		base = debian
	}

	node := "node.sh"
	nodePath := "scripts/" + node

	dir := c.Directory()

	nodeInstall, err := Scripts.ReadFile(nodePath)

	if err != nil {
		return nil, fmt.Errorf("unable to read embedded script %w", err)
	}

	golang := "golang.sh"
	golangPath := "scripts/" + golang

	golangInstall, err := Scripts.ReadFile(golangPath)

	if err != nil {
		return nil, fmt.Errorf("unable to read embedded script %w", err)
	}

	cloudKit := "cdk8s.sh"
	cloudKitPath := "scripts/" + cloudKit

	cloudKitInstall, err := Scripts.ReadFile(cloudKitPath)

	if err != nil {
		return nil, fmt.Errorf("unable to read embedded script %w", err)
	}

	dir = dir.WithNewFile("/"+node, string(nodeInstall))
	dir = dir.WithNewFile("/"+golang, string(golangInstall))
	dir = dir.WithNewFile("/"+cloudKit, string(cloudKitInstall))

	cdk := base.WithEntrypoint([]string{}).
		WithEnvVariable("VERSION", "18.14.2").
		WithEnvVariable("NODEGYPDEPENDENCIES", "true").
		WithMountedDirectory("/scripts", dir).
		WithExec(strings.Split("chmod +x "+"/scripts/node.sh", " ")).
		WithExec(strings.Split("/scripts/node.sh", " ")).
		WithEnvVariable("VERSION", "1.20.1").
		WithExec(strings.Split("chmod +x "+"/scripts/golang.sh", " ")).
		WithExec(strings.Split("/scripts/golang.sh", " ")).
		WithExec(strings.Split("chmod +x "+"/scripts/cdk8s.sh", " ")).
		WithExec(strings.Split("/scripts/cdk8s.sh", " "))

	return cdk, nil
}
