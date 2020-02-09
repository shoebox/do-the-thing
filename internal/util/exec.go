package util

import (
	"os/exec"
)

type Exec interface {
	Exec(dir *string, name string, args ...string) ([]byte, error)
}

type OsExec struct{}

func (e OsExec) Exec(workingDir *string, cmdName string, cmdArgs ...string) ([]byte, error) {
	exec := exec.Command(cmdName, cmdArgs...)
	if workingDir != nil {
		exec.Dir = *workingDir
	}
	return exec.Output()
}
