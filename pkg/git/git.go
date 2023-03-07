package git

import (
	"dagger.io/dagger"
)

type Opts struct {
	Version string
}

func Git(c *dagger.Client, opts Opts) (*dagger.Container, error) {

	return c.Container().From("bitnami/git:2.39.2-debian-11-r7").WithEntrypoint([]string{}), nil
}
