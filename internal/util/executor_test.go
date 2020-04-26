package util

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// setup:
	exec := NewExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// when:
	b, err := exec.CommandContext(ctx, "echo", "toto", "tata").Output()

	// then:
	assert.EqualValues(t, string(b), "toto tata\n")
	assert.NoError(t, err)
}

func TestTimeout(t *testing.T) {
	// setup:
	exec := NewExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	// when:
	err := exec.CommandContext(ctx, "sleep", "2").Run()

	// then:
	assert.EqualValues(t, err, context.DeadlineExceeded)
}

func TestSetEnv(t *testing.T) {
	// setup:
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ex := NewExecutor()

	// when:
	out, err := ex.CommandContext(ctx, "/bin/sh", "-c", "echo $TESTENVVAR").CombinedOutput()

	// then:
	assert.NoError(t, err)
	assert.EqualValues(t, out, "\n")

	// when:
	cmd := ex.CommandContext(ctx, "/bin/sh", "-c", "echo $TESTENVVAR")
	cmd.SetEnv([]string{"TESTENVVAR=xcode"})
	out, err = cmd.CombinedOutput()

	assert.NoError(t, err)
	assert.EqualValues(t, string(out), "xcode\n")
}

func TestStopBeforeStart(t *testing.T) {
	exec := NewExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	// when:
	cmd := exec.CommandContext(ctx, "sleep", "2")

	// no panic calling Stop before calling Run
	cmd.Stop()

	cmd.Run()

	// no panic calling Stop after command is done
	cmd.Stop()
}
