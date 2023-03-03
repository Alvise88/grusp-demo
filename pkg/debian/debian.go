package debian

import (
	"fmt"
	"runtime"

	"dagger.io/dagger"
	"github.com/alvise88/grusp-demo/pkg/util"
)

type DebianOpts struct {
	Platform dagger.Platform
	Version  string
	Packages []struct {
		Name    string
		Version string
	}
}

func Debian(client *dagger.Client, opts DebianOpts) (*dagger.Container, error) {
	version := opts.Version

	if version == "" {
		version = "11.6"
	}

	platform := opts.Platform

	if platform == "" {
		goos := runtime.GOOS
		goarch := runtime.GOARCH
		platform = dagger.Platform(fmt.Sprintf("%s/%s", goos, goarch))
	}

	debian := client.Container(dagger.ContainerOpts{Platform: platform}).
		From(fmt.Sprintf("debian:%s", version)).
		WithExec(util.ToCommand("apt-get update"))

	for _, pkg := range opts.Packages {
		debian = debian.WithExec(util.ToCommand(fmt.Sprintf("apt-get -y install --no-install-recommends %s", pkg.Name)))
	}

	return debian, nil
}
