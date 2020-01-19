package util

import (
	"io"

	"howett.net/plist"
)

func DecodeFile(reader io.ReadSeeker, message interface{}) error {
	decoder := plist.NewDecoder(reader)
	return decoder.Decode(message)
}
