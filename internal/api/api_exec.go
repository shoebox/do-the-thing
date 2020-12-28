package api

import (
	"context"
	"io"
)

// Executor command helper base interface
type Executor interface {
	// CommandContext allow to execute a command with Context
	CommandContext(ctx context.Context, cmd string, args ...string) Cmd
}

// Cmd is an interface that wrap the Cmd action from os/exec in a more friendly API
type Cmd interface {
	// Run runs the command to the completion.
	Run() error

	// CombinedOutput returns combined stdout and stder
	CombinedOutput() ([]byte, error)

	// Output runs the command and returns standard output, but not standard err
	Output() ([]byte, error)

	// SetDir set working directory
	SetDir(dir string)

	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(out io.Writer)

	// SetEnv allow to define environment values for the command
	SetEnv(env []string)

	// StdoutPipe and StderrPipe for getting the process outputs Stdout and Stderr as Readers
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)

	// Start and Wait are for running a process non-blocking
	Start() error
	Wait() error

	// Stops the command by sending SIGTERM. It is not guaranteed the
	// process will stop before this function returns. If the process is not
	// responding, an internal timer function will send a SIGKILL to force
	// terminate after 10 seconds.
	Stop()
}

