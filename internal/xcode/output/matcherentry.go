package output

import "regexp"

type matcherEntry struct {
	logfunc func(e LogEntry)
	regexp  *regexp.Regexp
}

func createMatcherEntry(f func(LogEntry), str string) matcherEntry {
	return matcherEntry{
		logfunc: f,
		regexp:  regexp.MustCompile(str),
	}
}

func (e *matcherEntry) Match(txt string) (bool, map[string]string) {
	match := e.regexp.MatchString(txt)
	res := make(map[string]string)

	if match {
		m := e.regexp.FindStringSubmatch(txt)
		for i, name := range e.regexp.SubexpNames() {
			if i > 0 && i <= len(m) {
				res[name] = m[i]
			}
		}
	}

	return match, res
}
