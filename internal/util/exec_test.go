package util

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var subject Exec

func TestMain(m *testing.M) {
	subject = NewCommandRunner()
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		// Normal test mode
		os.Exit(m.Run())

	case "stderrfail":
		fmt.Fprintf(os.Stderr, "some stderr text\n")
		os.Exit(1)

	case "echo":
		fmt.Println("Echo")
		iargs := []interface{}{}
		for _, s := range os.Args[1:] {
			iargs = append(iargs, s)
		}
		fmt.Println(" >>>> ", os.Args)
	}
}

func TestEchoWithContext(t *testing.T) {
	// setup:
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// when:
	b, err := subject.ContextExec(ctx, "echo", "toto", "tata")

	// then: Should echo thhe result
	assert.EqualValues(t, "toto tata\n", string(b))
	assert.NoError(t, err)
}

func TestEcho(t *testing.T) {
	// when:
	b, err := subject.Exec(nil, "echo", "toto", "tata")

	// then: Should echo thhe result
	assert.EqualValues(t, "toto tata\n", string(b))
	assert.NoError(t, err)
}

func TestEchoDir(t *testing.T) {
	// when:
	dir := "/path/to/dir"
	b, err := subject.Exec(&dir, "echo", "toto", "tata")

	// then: Should echo thhe result
	assert.EqualValues(t, "", string(b))

	// and: Expect an error
	assert.EqualError(t, err, "chdir /path/to/dir: no such file or directory")
}

func TestEchoFail(t *testing.T) {
	// setup:
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// when:
	b, err := subject.ContextExec(ctx, "stderrfail", "toto", "tata")

	// then: Should echo thhe result
	assert.Empty(t, b)
	assert.EqualError(t, err, "exec: \"stderrfail\": executable file not found in $PATH")
}
