package utiltest

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/stretchr/testify/mock"
)

type MockExecutorCmd struct {
	mock.Mock
}

func (cmd *MockExecutorCmd) SetDir(dir string) {
	cmd.Called(dir)
}

func (cmd *MockExecutorCmd) SetStdin(in io.Reader) {
	cmd.Called(in)
}

func (cmd *MockExecutorCmd) SetStdout(out io.Writer) {
	cmd.Called(out)
}

func (cmd *MockExecutorCmd) SetStderr(out io.Writer) {
	cmd.Called(out)
}

func (cmd *MockExecutorCmd) SetEnv(env []string) {
	cmd.Called(env)
}

func (cmd *MockExecutorCmd) StdoutPipe() (io.ReadCloser, error) {
	r := cmd.Called()
	return ioutil.NopCloser(bytes.NewBufferString(r.String(0))), r.Error(1)
}

func (cmd *MockExecutorCmd) StderrPipe() (io.ReadCloser, error) {
	r := cmd.Called()
	return ioutil.NopCloser(bytes.NewBufferString(r.String(0))), r.Error(1)
}

func (cmd *MockExecutorCmd) Start() error {
	r := cmd.Called()
	return r.Error(0)
}

func (cmd *MockExecutorCmd) Wait() error {
	r := cmd.Called()
	return r.Error(0)
}

func (cmd *MockExecutorCmd) Run() error {
	r := cmd.Called()
	return r.Error(0)
}

func (cmd *MockExecutorCmd) CombinedOutput() ([]byte, error) {
	r := cmd.Called()
	return []byte(r.String(0)), r.Error(1)
}

func (cmd *MockExecutorCmd) Output() ([]byte, error) {
	r := cmd.Called()
	return []byte(r.String(0)), r.Error(1)
}

// Stop is part of the Cmd interface.
func (cmd *MockExecutorCmd) Stop() {
	cmd.Called()
}
