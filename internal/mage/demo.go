package mage

import (
	"context"
	"os"

	"dagger.io/dagger"
	"github.com/alvise88/grusp-demo/internal/mage/util"
	"github.com/alvise88/grusp-demo/pkg/k8s"
	asbtractionutil "github.com/alvise88/grusp-demo/pkg/util"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

type Demo mg.Namespace

// build hello-grusp
func (t Demo) Build(ctx context.Context) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	repo := util.HelloRepository(c)

	_, err = c.Container().Build(repo).ExitCode(ctx)

	if err != nil {
		return err
	}

	return nil
}

// publish hello-grusp
func (t Demo) Publish(ctx context.Context) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	imageRepo := "alvisevitturi/hello-grusp:latest"

	repo := util.HelloRepository(c)

	hello := c.Container().Build(repo)

	_, err = hello.Publish(ctx, imageRepo)

	if err != nil {
		return err
	}

	return nil
}

// publish hello-grusp
func (t Demo) Deploy(ctx context.Context) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	cdkCode := c.Host().Directory("./infra", dagger.HostDirectoryOpts{
		Exclude: []string{"import/", "dist/"},
	})

	cdkCommand := []string{"cdk8s", "synth"}

	cdk, err := k8s.Cdk8s(c, k8s.Cdk8sOpts{
		Version: "2.1.148",
	})

	if err != nil {
		return err
	}

	kubeConfig := c.Host().Directory(os.ExpandEnv("${HOME}/.kube"))

	_, err = cdk.WithMountedDirectory("/opt/app", cdkCode).
		WithMountedDirectory("/root/.kube", kubeConfig).
		WithWorkdir("/opt/app").
		// WithExec(asbtractionutil.ToCommand("cdk8s --version")).
		WithExec(asbtractionutil.ToCommand("cdk8s import")).
		WithExec(cdkCommand).
		// WithExec(asbtractionutil.ToCommand("ls -la")).
		WithExec(asbtractionutil.ToCommand("kubectl apply -f ./dist/hello.k8s.yaml")).
		ExitCode(ctx)

	return err
}
