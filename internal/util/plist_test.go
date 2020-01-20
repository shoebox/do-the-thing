package util

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	BundleVersion string `plist:"CFBundleVersion"`
	Version       string `plist:"CFBundleShortVersionString"`
}

const (
	XCodeSample string = `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>CFBundleVersion</key>
			<string>1.2.3</string>
			<key>CFBundleShortVersionString</key>
			<string>1.2.3-snapshot</string>
		</dict>
	</plist>`
)

func TestDecoding(t *testing.T) {
	data := test{}

	t.Run("Should fail and report error in case of invalid payload", func(t *testing.T) {
		err := DecodeFile(bytes.NewReader([]byte("invalid")), data)
		fmt.Println("err ", err)

		assert.EqualError(t, err, "plist: type mismatch: tried to decode plist type `string' into value of type `util.test'")
	})

	t.Run("Should decode valid payload without issue", func(t *testing.T) {
		err := DecodeFile(bytes.NewReader([]byte(XCodeSample)), &data)

		assert.Nil(t, err)
		assert.EqualValues(t, data, test{BundleVersion: "1.2.3", Version: "1.2.3-snapshot"})
	})
}
