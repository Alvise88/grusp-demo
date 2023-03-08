package util

import (
	"dagger.io/dagger"
	"github.com/alvise88/grusp-demo/pkg/util"
)

// Repository with common set of exclude filters to speed up upload
func HelloRepository(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory("./hello", dagger.HostDirectoryOpts{
		Exclude: []string{
			// node
			"**/node_modules",

			// python
			"**/__pycache__",
			"**/.venv",
			"**/.mypy_cache",
			"**/.pytest_cache",
		},
	})
}

// Repository with common set of exclude filters to speed up upload
func InfraRepository(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory("./hello", dagger.HostDirectoryOpts{
		Exclude: []string{
			// node
			"**/node_modules",

			// python
			"**/__pycache__",
			"**/.venv",
			"**/.mypy_cache",
			"**/.pytest_cache",
		},
	})
}

// HelloGoCodeOnly is Repository, filtered to only contain Go code.
//
// NOTE: this function is a shared util ONLY because it's used both by the Engine
// and the Go SDK. Other languages shouldn't have a common helper.
func HelloGoCodeOnly(c *dagger.Client) *dagger.Directory {
	return c.Directory().WithDirectory("/", HelloRepository(c), dagger.DirectoryWithDirectoryOpts{
		Include: []string{
			// go source
			"**/*.go",

			// modules
			"**/go.mod",
			"**/go.sum",

			// embedded files
			"**/*.go.tmpl",
			"**/*.graphqls",
			"**/*.graphql",

			// custom embedded files
			"**/*.txt",
			"**/*.yaml",
			"**/*.sh",

			// misc
			".golangci.yml",
			"**/Dockerfile", // needed for shim TODO: just build shim directly
		},
	})
}

func GoBase(c *dagger.Client) *dagger.Container {
	repo := HelloGoCodeOnly(c)

	// Create a directory containing only `go.{mod,sum}` files.
	goMods := c.Directory()
	for _, f := range []string{"go.mod", "go.sum"} {
		goMods = goMods.WithFile(f, repo.File(f))
	}

	return c.Container().
		From("golang:1.20.1-alpine").
		// From("golang:1.20.0-alpine").
		// gcc is needed to run go test -race https://github.com/golang/go/issues/9918 (???)
		WithExec(util.ToCommand("apk add build-base")).
		WithEnvVariable("CGO_ENABLED", "0").
		// adding the git CLI to inject vcs info
		// into the go binaries
		WithExec([]string{"apk", "add", "git"}).
		WithWorkdir("/app").
		// run `go mod download` with only go.mod files (re-run only if mod files have changed)
		WithMountedDirectory("/app", goMods).
		WithExec([]string{"go", "mod", "download"}).
		// run `go build` with all source
		WithMountedDirectory("/app", repo).
		// include a cache for go build
		WithMountedCache("/root/.cache/go-build", c.CacheVolume("go-build"))
}
