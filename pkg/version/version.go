package version

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/alvise88/grusp-demo/pkg/git"
	"github.com/alvise88/grusp-demo/pkg/util"
)

type Opts struct {
	Base string

	Source *dagger.Directory
}

func Version(c *dagger.Client, opts Opts) (*dagger.Container, error) {
	script := "version.sh"
	scriptPath := "scripts/" + script

	dir := c.Directory()

	vsr, err := Scripts.ReadFile(scriptPath)

	if err != nil {
		return nil, fmt.Errorf("unable to read embedded script %w", err)
	}

	dir = dir.WithNewFile("/"+script, string(vsr))

	git, err := git.Git(c, git.Opts{})

	if err != nil {
		return nil, err
	}

	semVer := git.
		WithEnvVariable("BASE", opts.Base).
		WithMountedDirectory("/scripts", dir).
		WithMountedDirectory("/src", opts.Source).
		WithWorkdir("/src").
		WithEntrypoint([]string{}).
		WithExec(util.ToCommand("chmod +x " + "/scripts/version.sh")).
		WithExec(util.ToCommand("/scripts/version.sh"))

	return semVer, nil
}
