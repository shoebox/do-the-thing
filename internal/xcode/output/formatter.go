package output

import (
	"bufio"
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type formatter struct {
	m []matcherEntry
}

func NewFormatter(r reporter) formatter {
	return formatter{m: NewMatcher(r)}
}

func (frm formatter) Parse(r io.Reader, errType bool) {
	entry := LogEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		// log.Debug().Msg(txt)
		var found bool
		for _, matcher := range frm.m {
			match, m := matcher.Match(strings.ReplaceAll(txt, `\ `, ""))
			if match {
				err := mapstructure.Decode(m, &entry)

				if err == nil {
					matcher.logfunc(entry)
					break
				}

				found = true
			}
		}

		if !found && errType {
			log.Error().Msg(txt)
		}
	}
}
