package keychain

import (
	"dothething/internal/api"
	"dothething/internal/utiltest"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type keychainSearchTestSuite struct {
	suite.Suite
	API     *api.API
	exec    *utiltest.MockExecutor
	subject keychain
}

func TestKeychainSearchListSuite(t *testing.T) {
	suite.Run(t, new(keychainSearchTestSuite))
}

func (s *keychainSearchTestSuite) TestGetSearchListInputOutputAndParsing() {
	// setup:
	cases := []struct {
		expected []string
		output   string
		err      error
	}{
		{
			output:   `    "/Users/user.name/Library/Keychains/login.keychain-db"`,
			expected: []string{"/Users/user.name/Library/Keychains/login.keychain-db"},
			err:      nil,
		},
		{
			output: `    "/Users/user.name/Library/Keychains/login.keychain-db"
    "/Users/user.name/Library/Keychains/test2.keychain-db"`,
			expected: []string{
				"/Users/user.name/Library/Keychains/login.keychain-db",
				"/Users/user.name/Library/Keychains/test2.keychain-db",
			},
			err: nil,
		},
		{
			output: "",
			err:    errors.New("toto"),
		},
	}

	for _, c := range cases {
		// when:
		r := parseSearchList([]byte(c.output))

		// and:
		s.Assert().EqualValues(c.expected, r)
	}
}

func (s *keychainSearchTestSuite) TestParseSearchListLine() {
	cases := []struct {
		line     string
		expected string
	}{
		{
			line:     `    "/Users/user.name/Library/Keychains/login.keychain-db"`,
			expected: `/Users/user.name/Library/Keychains/login.keychain-db`,
		},
		{
			line:     `    "/private/var/folders/m1/s05mrl8s4fbfw8zf4hw04tjwh740pl/T/do-the-thing-856917238/do-the-thing.keychain"`,
			expected: "/private/var/folders/m1/s05mrl8s4fbfw8zf4hw04tjwh740pl/T/do-the-thing-856917238/do-the-thing.keychain",
		},
		{
			line:     "/path/to",
			expected: "",
		},
		{
			line:     `    /path/to`,
			expected: "",
		},
	}

	for _, c := range cases {
		s.Run(c.line, func() {
			// when:
			res := isSearchListEntry(c.line)

			// then:
			s.Assert().EqualValues(res, c.expected)
		})
	}
}
