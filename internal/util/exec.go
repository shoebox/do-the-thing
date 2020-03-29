package util

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Exec interface {
	Exec(dir *string, name string, args ...string) ([]byte, error)
	ContextExec(ctx context.Context, name string, args ...string) ([]byte, error)
}

type OsExec struct{}

// ContextExec execute the program with the provided arguments and context
func (e OsExec) ContextExec(ctx context.Context, name string, args ...string) ([]byte, error) {
	log.Debug().Str("name", name).Strs("Args", args).Msg("Execute with context")
	//Msg("Execute with context")
	stdout := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = stdout

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return stdout.Bytes(), nil
}

// Exec execute the program with the provided arguments
func (e OsExec) Exec(workingDir *string, cmdName string, cmdArgs ...string) ([]byte, error) {
	exec := exec.Command(cmdName, cmdArgs...)
	if workingDir != nil {
		exec.Dir = *workingDir
	}
	return exec.Output()
}
