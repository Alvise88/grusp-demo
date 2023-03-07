package util

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

func ToCommand(cmd string) []string {
	return strings.Split(cmd, " ")
}

func WithSetHostVar(ctx context.Context, h *dagger.Host, varName string) *dagger.HostVariable {
	hv := h.EnvVariable(varName)
	if val, _ := hv.Secret().Plaintext(ctx); val == "" {
		fmt.Fprintf(os.Stderr, "env var %s must be set", varName)
		os.Exit(1)
	}
	return hv
}
