package xcode

import (
	"fmt"
	"regexp"
	"strconv"
)

type errorCode int

const (
	// USAGE The command was used incorrectly, e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
	USAGE errorCode = 64 + iota

	// DATAERR The input data was incorrect in some way.  This should only be used for user's data and not system files.
	DATAERR

	// NOINPUT An input file (not a system file) did not exist or was not readable.  This could also include errors like ``No message'' to a mailer (if it cared to catch it).
	NOINPUT

	// NOUSER The user specified did not exist.  This might be used for mail addresses or remote logins.
	NOUSER

	// NOHOST The host specified did not exist.  This is used in mail addresses or network requests.
	NOHOST

	// UNAVAILABLE A service is unavailable.  This can occur if a support program or file does not exist.  This can also be used as a catchall message when something you wanted to do doesn't work, but you don't know why.
	UNAVAILABLE

	// SOFTWARE An internal software error has been detected.  This should be limited to non-operating system related errors as possible.
	SOFTWARE

	// OSERR An operating system error has been detected.  This is intended to be used for such things as ``cannot fork'', ``cannot create pipe'', or the like.  It includes things like getuid returning a user that does not exist in the passwd file.
	OSERR

	// OSFILE Some system file (e.g., /etc/passwd, /var/run/utmp, etc.) does not exist, cannot be opened, or has some sort of error (e.g., syntax error).
	OSFILE

	// CANTCREAT A (user specified) output file cannot be created.
	CANTCREAT

	// IOERR An error occurred while doing I/O on some file.
	IOERR

	// TEMPFAIL Temporary failure, indicating something that is not really an error.  In sendmail, this means that a mailer (e.g.) could not create a connection, and the request should be reattempted later.
	TEMPFAIL

	// PROTOCOL The remote system returned something that was ``not possible'' during a protocol exchange.
	PROTOCOL

	// NOPERM You did not have sufficient permission to perform the operation.  This is not intended for file system problems, which should use EX_NOINPUT or EX_CANTCREAT, but rather for higher level permissions.
	NOPERM

	// CONFIG Something was found in an unconfigured or misconfigured state.
	CONFIG
)

// NewError create a new wrapped xcodebuild error from the provided error code
func NewError(code int) XCodebuildError {
	return XCodebuildError{code: errorCode(code), errorCode: code}
}

// XCodebuildError wrap the error code in a more readable struct
type XCodebuildError struct {
	errorCode int
	code      errorCode
}

// Error implement the error interface and return a readable description of the error
func (e XCodebuildError) Error() string {
	var msg string
	switch e.code {
	case USAGE:
		msg = "Command configuration error"
	case DATAERR:
		msg = "The input data was incorrect in some way."
	case NOINPUT:
		msg = "An input file (not a system file) did not exist or was not readable."
	case NOUSER:
		msg = "The user specified did not exist."
	case NOHOST:
		msg = "The host specified did not exist."
	case UNAVAILABLE:
		msg = "A service is unavailable."
	case SOFTWARE:
		msg = "An internal software error has been detected."
	case OSERR:
		msg = "An operating system error has been detected."
	case OSFILE:
		msg = "Some system file (e.g., /etc/passwd, /var/run/utmp, etc.) does not exist"
	case CANTCREAT:
		msg = "A (user specified) output file cannot be created."
	case IOERR:
		msg = "An error occurred while doing I/O on some file."
	case TEMPFAIL:
		msg = "Temporary failure, indicating something that is not really an error"
	case PROTOCOL:
		msg = "The remote system returned something that was ``not possible'' during a protocol exchange."
	case NOPERM:
		msg = "You did not have sufficient permission to perform the operation."
	case CONFIG:
		msg = "Something was found in an unconfigured or misconfigured state."
	default:
		msg = "Unknown error"
	}

	return fmt.Sprintf("Error %v - %v", e.code, msg)
}

var errorRegExp = regexp.MustCompile(`^exit\sstatus\s(?P<Code>[0-9]+)`)

func ParseXCodeBuildError(xerr error) error {
	if xerr == nil {
		return nil
	}

	return parseError(xerr.Error())
}

func parseError(txt string) error {
	var index = -1
	if errorRegExp.MatchString(txt) {
		m := errorRegExp.FindStringSubmatch(txt)
		index = parseErrorIntIndex(m[1])
	}

	return NewError(index)
}

func parseErrorIntIndex(txt string) int {
	index, err := strconv.Atoi(txt)
	if err != nil {
		return -1
	}
	return index
}
