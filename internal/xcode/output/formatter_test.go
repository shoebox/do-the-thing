package output

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func TestMain(t *testing.M) {
	// assert.NoError(t, err)
}
*/

func TestParsing(t *testing.T) {
	fmt.Println("")
	data, err := os.Open("testdata/output.txt")
	assert.NoError(t, err)
	Parse(data)

	//
	// fmt.Println(data)
}
