package xcode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseErrorIndex(t *testing.T) {
	// setup:
	cases := []struct {
		txt      string
		expected int
	}{
		{txt: "42", expected: 42},
		{txt: "40", expected: 40},
		{txt: "", expected: -1},
	}

	// when:
	for i, tc := range cases {
		t.Run(fmt.Sprintf("Case %v", i), func(t *testing.T) {
			v := parseErrorIntIndex(tc.txt)

			// then:
			assert.EqualValues(t, tc.expected, v)
		})
	}
}

func TestInvalidErrorParsing(t *testing.T) {
	// when:
	e := parseError("Invalid error 8")

	// then:
	assert.EqualError(t, e, "Error -1 - Unknown error")
}

func TestErrorParsing(t *testing.T) {
	// when:
	e := parseError("exit status 70")

	// then:
	assert.EqualError(t, e, "Error 70 - An internal software error has been detected.")
}
