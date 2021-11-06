package keychain

import (
	"context"
	"dothething/internal/api"
)

func (k keychain) securityCmd(ctx context.Context, action string, args []string) api.Cmd {
	return k.API.Exec.CommandContext(
		ctx,
		SecurityUtil,
		append([]string{action}, args...)...,
	)
}
