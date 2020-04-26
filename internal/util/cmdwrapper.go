package util

import (
	"io"
	"os/exec"
	"syscall"
	"time"
)

// Wraps exec.Cmd so we can capture errors.
type cmdWrapper exec.Cmd

var _ Cmd = &cmdWrapper{}

// SetDir is part of the Cmd interface
func (cmd *cmdWrapper) SetDir(dir string) {
	cmd.Dir = dir
}

// SetStdin is part of the Cmd interface
func (cmd *cmdWrapper) SetStdin(in io.Reader) {
	cmd.Stdin = in
}

// SetStdout is part of the Cmd interface
func (cmd *cmdWrapper) SetStdout(out io.Writer) {
	cmd.Stdout = out
}

// SetStderr is part of the Cmd interface
func (cmd *cmdWrapper) SetStderr(out io.Writer) {
	cmd.Stderr = out
}

// SetEnv is part of the Cmd interface
func (cmd *cmdWrapper) SetEnv(env []string) {
	cmd.Env = env
}

// StdoutPipe is part of the Cmd interface
func (cmd *cmdWrapper) StdoutPipe() (io.ReadCloser, error) {
	r, err := (*exec.Cmd)(cmd).StdoutPipe()
	return r, handleError(err)
}

// StderrPipe is part of the Cmd interface
func (cmd *cmdWrapper) StderrPipe() (io.ReadCloser, error) {
	r, err := (*exec.Cmd)(cmd).StderrPipe()
	return r, handleError(err)
}

// Start is part of the Cmd interface
func (cmd *cmdWrapper) Start() error {
	err := (*exec.Cmd)(cmd).Start()
	return handleError(err)
}

// Wait is part of the Cmd interface
func (cmd *cmdWrapper) Wait() error {
	err := (*exec.Cmd)(cmd).Wait()
	return handleError(err)
}

// Run is part of the Cmd interface
func (cmd *cmdWrapper) Run() error {
	err := (*exec.Cmd)(cmd).Run()
	return handleError(err)
}

// CombinedOutput is part of the Cmd
func (cmd *cmdWrapper) CombinedOutput() ([]byte, error) {
	out, err := (*exec.Cmd)(cmd).CombinedOutput()
	return out, handleError(err)
}

// Output is part of the Cmd interface.
func (cmd *cmdWrapper) Output() ([]byte, error) {
	out, err := (*exec.Cmd)(cmd).Output()
	return out, handleError(err)
}

// Stop is part of the Cmd interface.
func (cmd *cmdWrapper) Stop() {
	c := (*exec.Cmd)(cmd)

	if c.Process == nil {
		return
	}

	c.Process.Signal(syscall.SIGTERM)

	time.AfterFunc(10*time.Second, func() {
		if !c.ProcessState.Exited() {
			c.Process.Signal(syscall.SIGKILL)
		}
	})
}
