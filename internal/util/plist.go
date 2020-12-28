package util

import (
	"io"

	"howett.net/plist"
)

type DecodingError struct {
}

func (e DecodingError) Error() string {
	return "Failed to decode plist payload"
}

func DecodeFile(reader io.ReadSeeker, message interface{}) error {
	decoder := plist.NewDecoder(reader)
	if err := decoder.Decode(message); err != nil {
		return DecodingError{}
	}

	return nil
}
