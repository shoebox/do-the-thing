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
	fmt.Println("TestMain")
	subject = NewCommandRunner()
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		// Normal test mode
		os.Exit(m.Run())

	case "stderrfail":
		fmt.Fprintf(os.Stderr, "some stderr text\n")
		os.Exit(1)

	case "echo":
		iargs := []interface{}{}
		for _, s := range os.Args[1:] {
			iargs = append(iargs, s)
		}
		fmt.Println(iargs...)
	}
}

func TestEcho(t *testing.T) {
	// setup:
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// when:
	b, err := subject.ContextExec(ctx, "echo", "toto", "tata")

	// then: Should echo thhe result
	assert.EqualValues(t, "toto tata\n", string(b))
	assert.NoError(t, err)
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
