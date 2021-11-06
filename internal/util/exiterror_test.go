package util

/*
func TestExitStatus(t *testing.T) {
	ctx := context.Background()

	// when:
	_, err := NewExecutor().CommandContext(ctx, "fake-command").Output()

	// then:
	assert.EqualError(t, err, ErrExecutableNotFound.Error())
}

func TestExitStatusCode(t *testing.T) {
	ctx := context.Background()

	// when:
	_, err := NewExecutor().CommandContext(ctx, "cat", "/toto/tutu/fake.txt").Output()

	// then:
	assert.IsType(t, new(ExitErrorWrapper), err)

	// and:
	v := err.(*ExitErrorWrapper)
	assert.EqualValues(t, 1, v.ExitError.ExitCode())
	assert.EqualValues(t, 1, v.ExitStatus())
	assert.EqualValues(t, "exit status 1", v.Error())
}
*/
