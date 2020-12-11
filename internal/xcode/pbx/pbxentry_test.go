package pbx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var entry Entry
var ref1 Entry
var ref2 Entry
var raw PBXProjRaw

func TestMain(m *testing.M) {
	entry = Entry{Name: "Test name"}
	ref1 = Entry{Name: "ref1"}
	ref2 = Entry{Name: "ref2"}

	raw = PBXProjRaw{
		Objects: map[string]Entry{
			"test": entry,
			"ref1": ref1,
			"ref2": ref2,
		},
	}

	os.Exit(m.Run())
}

func TestGetRoot(t *testing.T) {
	// when:
	r := raw.GetRoot()

	// then:
	assert.ObjectsAreEqual(entry, r)
}

func TestGetRef(t *testing.T) {
	// setup:
	ref := Ref("test")

	// when:
	e := ref.Get(raw)

	// then:
	assert.ObjectsAreEqual(entry, e)
}
