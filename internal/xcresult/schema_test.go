package xcresult

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIssues(t *testing.T) {
	// setup:
	b, err := ioutil.ReadFile("testdata/xcresult3_summary.json")
	assert.NoError(t, err)

	// when:
	res, err := ParseIssues(b)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(res))

	// then:
	assert.EqualValues(t,
		res[0].Message,
		IssueType{
			Type: SupertypeClass{
				Name: "String",
			},
			Value: "XCTAssertEqual failed: (\"100\") is not equal to (\"101\") - fake test",
		})
}
