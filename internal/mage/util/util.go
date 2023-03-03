package util

import (
	"dagger.io/dagger"
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
