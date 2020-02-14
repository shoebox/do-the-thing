package signature

import (
	"dothething/internal/utiltest"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockExec *utiltest.MockExec

func setup() {
	mockExec = new(utiltest.MockExec)
}

func TestProvisioningProfileDecoding(t *testing.T) {
	path := "test/path/test.mobileprovision"

	t.Run("Decoding erors should be raised", func(t *testing.T) {
		setup()

		// setup:
		mockExec.
			On("Exec", Security, []string{Cms, ArgDecodeCMS, ArgInlineFile, path}).
			Return(nil, errors.New("error"))

		// when:
		pp, err := DecodeProvisioningProfile(path, mockExec)

		// then:
		assert.EqualValues(t, err, ErrorFailedToDecode)
		assert.Nil(t, pp)
	})

	t.Run("Decoding erors should be raised", func(t *testing.T) {
		setup()

		// setup:
		mockExec.
			On("Exec", Security, []string{Cms, ArgDecodeCMS, ArgInlineFile, path}).
			Return("invalid", nil)

		// when:
		pp, err := DecodeProvisioningProfile(path, mockExec)

		// then:
		assert.EqualValues(t, err, ErrorFailedToDecode)
		assert.Nil(t, pp)
	})
}

func TestParseRawX509Certificates(t *testing.T) {
	t.Run("Should throw an error if failed to parse the provisiong public key", func(t *testing.T) {
		// setup:
		data := [][]byte{[]byte("Hello world")}

		// when:
		res, err := parseRawX509Certificates(data)

		// then:
		assert.EqualValues(t, ErrorParsingPublicKey, err)
		assert.Nil(t, res)
	})
}

func TestFileDecodingExecution(t *testing.T) {
	t.Run("Should execute the command", func(t *testing.T) {
		setup()
		path := "test/path/test.mobileprovision"

		// setup:
		mockExec.
			On("Exec", Security, []string{Cms, ArgDecodeCMS, ArgInlineFile, path}).
			Return("hello world", nil)

		// when:
		b, err := decodeProvisioning(path, mockExec)

		// then:
		assert.NoError(t, err)
		assert.EqualValues(t, []byte("hello world"), b)

		// and:
		mockExec.AssertExpectations(t)
	})

	t.Run("Should handle command errorr", func(t *testing.T) {
		setup()
		path := "test/path/test.mobileprovision"

		// setup:
		mockExec.
			On("Exec", Security, []string{Cms, ArgDecodeCMS, ArgInlineFile, path}).
			Return(nil, errors.New("error"))

		// when:
		b, err := decodeProvisioning(path, mockExec)

		// then:
		assert.EqualValues(t, ErrorFailedToDecode, err)
		assert.Nil(t, b)

		// and:
		mockExec.AssertExpectations(t)
	})
}
