package output

import (
	"bufio"
	"fmt"
	"io"
	"strings"

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
		fmt.Println(txt)
		for _, matcher := range f.m {
			b, m := matcher.Match(strings.ReplaceAll(txt, `\ `, ""))
			if b {
				if err := mapstructure.Decode(m, &entry); err == nil {
					matcher.logfunc(entry)
				}
				break
			}
		}
	}
}
