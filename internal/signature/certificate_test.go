package signature

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoding(t *testing.T) {
	t.Run("Valid decoding should succeed", func(t *testing.T) {
		// setup:
		data, err := os.Open("../../assets/Certificate.p12")
		assert.NoError(t, err)
		assert.NotNil(t, data)

		// when:
		b, err := DecodeCertificate(data, "p4ssword")

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
		data, err := DecodeCertificate(reader, "pass")

		// then:
		assert.NotNil(t, err)
		assert.EqualError(t, err, ErrorFailedToReadFile.Error())

		// and: nothing should be returned
		assert.Nil(t, data)
	})

	t.Run("Decoding error should be handled", func(t *testing.T) {
		// setup:
		reader := strings.NewReader("invalid")

		// when:
		data, err := DecodeCertificate(reader, "pass")

		// then:
		assert.NotNil(t, err)
		assert.EqualError(t, err, ErrorFailedToDecryptPEM.Error())

		// and: nothing should be returned
		assert.Nil(t, data)
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
