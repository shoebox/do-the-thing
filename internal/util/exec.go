package util

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type Exec interface {
	Exec(dir *string, name string, args ...string) ([]byte, error)
	ContextExec(ctx context.Context, name string, args ...string) ([]byte, error)
}

type OsExec struct{}

func (e OsExec) ContextExec(ctx context.Context, name string, args ...string) ([]byte, error) {
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

func (e OsExec) Exec(workingDir *string, cmdName string, cmdArgs ...string) ([]byte, error) {
	fmt.Println(" >>>> ", workingDir, cmdName, cmdArgs)
	exec := exec.Command(cmdName, cmdArgs...)
	if workingDir != nil {
		exec.Dir = *workingDir
	}
	return exec.Output()
}
