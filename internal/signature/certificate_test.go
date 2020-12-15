package signature

import (
	"dothething/internal/utiltest"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fs *utiltest.MockFileService
var cs certService

func Main(m *testing.M) {
	fs = new(utiltest.MockFileService)
	// cs = certService{fs: fs}

	os.Exit(m.Run())
}

func TestIsCertFile(t *testing.T) {
	// setup:
	cases := []struct {
		fi    utiltest.MockFileInfo
		name  string
		valid bool
	}{
		{
			fi:    utiltest.NewMockFileInfo(0, false, ""),
			name:  "Invalid file mode",
			valid: false,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeAppend, true, "toto.p12"),
			name:  "Should not be a directory",
			valid: false,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeAppend, false, "toto.p12"),
			name:  "Valid mode",
			valid: true,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeAppend, false, "toto.prov"),
			name:  "Should have the right extension",
			valid: false,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeIrregular, false, "toto.p12"),
			name:  "Should be regular mode",
			valid: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// when:
			res := isCertificateFile(c.fi)

			// then:
			assert.EqualValues(t, c.valid, res)
		})
	}
}

func TestDecoding(t *testing.T) {
	t.Run("Valid decoding should succeed", func(t *testing.T) {
		// setup:
		data, err := os.Open("../../assets/Certificate.p12")
		assert.NoError(t, err)
		assert.NotNil(t, data)

		// when:
		b, err := cs.DecodeCertificate(data, "p4ssword")

		// then:
		assert.NoError(t, err)

		// and:
		assert.Equal(t,
			"CN=iPhone Distribution: Dummy Name Ltd (12345ABCDE),OU=SELFSIGNED,C=GB",
			b.Issuer.String())
	})

	t.Run("Reading failure should be reported", func(t *testing.T) {
		// setup:
		reader := errReader(0)

		// when:
		data, err := cs.DecodeCertificate(reader, "pass")

		// then:
		assert.NotNil(t, err)
		assert.EqualError(t, err, ErrorFailedToReadFile.Error())

		// and: nothing should be returned
		assert.Empty(t, data)
	})

	t.Run("Decoding error should be handled", func(t *testing.T) {
		// setup:
		reader := strings.NewReader("invalid")

		// when:
		data, err := cs.DecodeCertificate(reader, "pass")

		// then:
		assert.NotNil(t, err)
		assert.EqualError(t, err, ErrorFailedToDecryptPEM.Error())

		// and: nothing should be returned
		assert.Empty(t, data)
	})
}

func TestFailureToReadFileShouldBeReporter(t *testing.T) {
	t.Run("Should report error when failing to read file", func(t *testing.T) {
		// setup:
		reader := errReader(0)

		// when: Trying to decode an empty file
		contents, err := readFile(reader)

		// then: An error should be fired
		assert.NotNil(t, err)
		assert.EqualError(t, err, ErrorFailedToReadFile.Error())

		// and: nothing should be returned
		assert.Nil(t, contents)
	})
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
