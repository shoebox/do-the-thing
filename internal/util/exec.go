package util

import "os/exec"

type Exec interface {
	Exec(name string, args ...string) ([]byte, error)
}

type OsExec struct{}

func (e OsExec) Exec(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}
