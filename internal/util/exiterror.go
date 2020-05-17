package util

import (
	"os/exec"
	"syscall"
)

// ExitErrorWrapper is an implementation of ExitError in terms of os/exec ExitError.
// Note: standard exec.ExitError is type *os.ProcessState, which already implements Exited().
type ExitErrorWrapper struct {
	*exec.ExitError
}

// ErrExecutableNotFound is returned if the executable is not found.
var ErrExecutableNotFound = exec.ErrNotFound

func handleError(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *exec.ExitError:
		return &ExitErrorWrapper{e}
	case *exec.Error:
		if e.Err == exec.ErrNotFound {
			return ErrExecutableNotFound
		}
	}

	return err
}

// ExitStatus is part of the ExitError interface.
func (w ExitErrorWrapper) ExitStatus() int {
	ws, ok := w.Sys().(syscall.WaitStatus)
	if !ok {
		panic("can't call ExitStatus() on a non-WaitStatus exitErrorWrapper")
	}
	return ws.ExitStatus()
}

// CodeExitError is an implementation of ExitError consisting of an error object
// and an exit code (the upper bits of os.exec.ExitStatus).
type CodeExitError struct {
	Err  error
	Code int
}

// ExitError is an interface that presents an API similar to os.ProcessState, which is
// what ExitError from os/exec is. This is designed to make testing a bit easier and
// probably loses some of the cross-platform properties of the underlying library.
type ExitError interface {
	String() string
	Error() string
	Exited() bool
	ExitStatus() int
}

var _ ExitError = CodeExitError{}

func (e CodeExitError) Error() string {
	return e.Err.Error()
}

func (e CodeExitError) String() string {
	return e.Err.Error()
}

// Exited is to check if the process has finished
func (e CodeExitError) Exited() bool {
	return true
}

// ExitStatus is for checking the error code
func (e CodeExitError) ExitStatus() int {
	return e.Code
}
