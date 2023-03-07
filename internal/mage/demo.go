package mage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"dagger.io/dagger"
	"github.com/alvise88/grusp-demo/internal/mage/util"
	"github.com/alvise88/grusp-demo/pkg/k8s"
	asbtractionutil "github.com/alvise88/grusp-demo/pkg/util"
	"github.com/alvise88/grusp-demo/pkg/version"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

type Demo mg.Namespace

var helloRepo = "alvisevitturi/hello-grusp"
var installerRepo = "alvisevitturi/hello-grusp-installer"

func dirty(ctx context.Context, c *dagger.Client) bool {
	clean, err := c.Host().EnvVariable("CI").Value(ctx)

	if err != nil {
		clean = "false"
	}

	cleanBool, err := strconv.ParseBool(clean)

	if err != nil {
		cleanBool = false
	}

	return !cleanBool
}

func vsr(ctx context.Context, c *dagger.Client) (string, error) {
	semVer, err := version.Version(c, version.Opts{
		Base:  "1.0",
		Dirty: dirty(ctx, c),
		Source: c.Host().Directory(".", dagger.HostDirectoryOpts{
			Include: []string{
				".git",
			},
		}),
	})

	if err != nil {
		return "", err
	}

	return semVer.Stdout(ctx)
}

func hello(c *dagger.Client) *dagger.Container {
	repo := util.HelloRepository(c)

	grusp := c.Container().Build(repo)

	return grusp
}

func installer(c *dagger.Client) (*dagger.Container, error) {
	cdkCode := c.Host().Directory("./infra", dagger.HostDirectoryOpts{
		Exclude: []string{"import/", "dist/"},
	})

	cdk, err := k8s.Cdk8s(c, k8s.Cdk8sOpts{
		Version: "2.1.148",
	})

	if err != nil {
		return nil, err
	}

	install := cdk.WithDirectory("/opt/app", cdkCode).
		WithWorkdir("/opt/app").
		WithExec(asbtractionutil.ToCommand("cdk8s import"))

	return install, nil
}

// build hello-grusp
func (t Demo) Version(ctx context.Context) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(io.Discard))
	if err != nil {
		return err
	}
	defer c.Close()

	vsr, err := vsr(ctx, c)

	if err != nil {
		return err
	}

	// fmt.Printf("version: %s", vsr)
	fmt.Print(vsr)

	return nil
}

// build hello-grusp
func (t Demo) Build(ctx context.Context) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	grusp := hello(c)

	_, err = grusp.ExitCode(ctx)

	return err
}

// publish hello-grusp
func (t Demo) Publish(ctx context.Context) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	grusp := hello(c)

	vsr, err := vsr(ctx, c)

	if err != nil {
		return err
	}

	_, err = grusp.Publish(ctx, fmt.Sprintf("%s:%s", helloRepo, vsr))

	if err != nil {
		return err
	}

	_, err = grusp.Publish(ctx, fmt.Sprintf("%s:%s", helloRepo, "latest"))

	if err != nil {
		return err
	}

	install, err := installer(c)

	if err != nil {
		return err
	}

	_, err = install.Publish(ctx, fmt.Sprintf("%s:%s", installerRepo, vsr))

	if err != nil {
		return err
	}

	_, err = install.Publish(ctx, fmt.Sprintf("%s:%s", installerRepo, "latest"))

	return err
}

// publish hello-grusp
func (t Demo) Deploy(ctx context.Context) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	vsr, err := vsr(ctx, c)

	if err != nil {
		return err
	}

	cdkCommand := []string{"cdk8s", "synth"}

	kubeConfig := c.Host().Directory(os.ExpandEnv("${HOME}/.kube"))

	varReplicas := c.Host().EnvVariable("HELLO_REPLICAS")

	replicas, err := varReplicas.Value(ctx)

	if err != nil {
		replicas = "1"
	}

	_, err = c.Container().From(fmt.Sprintf("%s:%s", installerRepo, vsr)).
		WithMountedDirectory("/root/.kube", kubeConfig).
		WithEnvVariable("HELLO_VERSION", vsr).
		WithEnvVariable("HELLO_REPLICAS", replicas).
		WithExec(cdkCommand).
		WithExec(asbtractionutil.ToCommand("kubectl apply -f ./dist/hello.k8s.yaml")).
		ExitCode(ctx)

	return err
}
