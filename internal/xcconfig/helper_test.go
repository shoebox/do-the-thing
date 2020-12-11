package xcconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var subject helper

func TestMain(m *testing.M) {
	subject = helper{entries: map[string]entry{}}

	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {

	// when:
	b, err := subject.Add("Foo", "value", EntryConfig{})

	// then:
	assert.True(t, b)
	assert.NoError(t, err)

	// when:
	b, err = subject.Add("Foo", "value", EntryConfig{})

	// and:
	assert.False(t, b)
	assert.EqualError(t, err, "Key 'Foo' is already being used")
}

func TestAddAndGenerate(t *testing.T) {
	// setup:
	subject = helper{entries: map[string]entry{}}
	b, err := subject.Add("CODE_SIGN_STYLE", "Manual", EntryConfig{Config: "BookStore"})
	assert.True(t, b)
	assert.NoError(t, err)

	// when:
	res := subject.Generate()

	// then:
	assert.EqualValues(t, "CODE_SIGN_STYLE[config=BookStore]=Manual", res)
}

func TestGenerate(t *testing.T) {
	cases := []struct {
		Name     string
		Case     map[string]entry
		Expected string
	}{
		{
			Name:     "basic",
			Case:     map[string]entry{"FOO": {Value: "dummy"}},
			Expected: "FOO=dummy",
		},
		{
			Name:     "arch constraint",
			Case:     map[string]entry{"FOO": {Value: "dummy", Config: EntryConfig{ARCH: "archname"}}},
			Expected: "FOO[arch=archname]=dummy",
		},
		{
			Name: "Single",
			Case: map[string]entry{
				"FOO": {
					Value:  "dummy",
					Config: EntryConfig{ARCH: "archname", Config: "configname"},
				},
			},
			Expected: "FOO[arch=archname][config=configname]=dummy",
		},
		{
			Name: "Multiple",
			Case: map[string]entry{
				"FOO1": {
					Value:  "dummy",
					Config: EntryConfig{ARCH: "archname", Config: "configname"},
				},
				"FOO2": {
					Value:  "dummy2",
					Config: EntryConfig{},
				},
			},
			Expected: "FOO1[arch=archname][config=configname]=dummy\nFOO2=dummy2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// setup:
			subject.entries = tc.Case

			// when:
			res := subject.Generate()

			assert.EqualValues(t, tc.Expected, res)
		})
	}
}

func TestFormatSDK(t *testing.T) {
	cases := []struct {
		Name     string
		Case     SDKConfig
		Expected string
	}{
		{Name: "Empty", Case: SDKConfig{}, Expected: ""},
		{
			Name:     "SDK Name only",
			Case:     SDKConfig{Name: "macosx"},
			Expected: "[sdk=macosx]",
		},
		{
			Name:     "SDK Name and version",
			Case:     SDKConfig{Name: "macosx", Version: "10.1"},
			Expected: "[sdk=macosx10.1]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// when:
			res := subject.formatSDKEntry(tc.Case)

			// then:
			assert.EqualValues(t, tc.Expected, res)
		})
	}
}

func TestFormatKeyValue(t *testing.T) {
	// setup:
	cases := []struct {
		Name     string
		Case     EntryConfig
		Expected string
	}{
		{Name: "No arch", Case: EntryConfig{}, Expected: ""},
		{Name: "i386", Case: EntryConfig{ARCH: "i386"}, Expected: "[arch=i386]"},
		{Name: "arm*", Case: EntryConfig{ARCH: "arm*"}, Expected: "[arch=arm*]"},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// when:
			res := subject.formatKeyValue("arch", tc.Case.ARCH)

			// then:
			assert.EqualValues(t, tc.Expected, res)
		})
	}
}

/*
func TestFormat(t *testing.T) {
	// setup:
	subject.Add("FOO", "dummy", EntryConfig{SDK: SDKConfig{Name: "iphoneos", Version: "1.2.3"}})

	// when:
	res := subject.Generate()

	// then:
	assert.EqualValues(t, "FOO[sdk=iphoneos1.2.3] = dummy", res)
}
*/
