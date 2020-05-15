package output

import (
	"bufio"
	"io"

	"github.com/mitchellh/mapstructure"
)

type formatter struct {
	m []matcherEntry
}

func NewFormatter(r reporter) formatter {
	return formatter{m: NewMatcher(r)}
}

func (f formatter) Parse(r io.Reader) {
	entry := LogEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		for _, matcher := range f.m {
			b, m := matcher.Match(txt)
			if b {
				mapstructure.Decode(m, &entry)
				matcher.logfunc(entry)
				break
			}
		}
	}
}
