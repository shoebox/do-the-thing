package output

import (
	"bufio"
	"io"
	"reflect"
	"regexp"

	"github.com/mitchellh/mapstructure"
)

func createEntry(f func(LogEntry), str string) entry {
	return entry{
		logfunc: f,
		regexp:  regexp.MustCompile(str),
	}
}

type entry struct {
	logfunc func(e LogEntry)
	regexp  *regexp.Regexp
}

type LogEntry struct {
	Entry             entry
	Arg               string
	BuildArch         string
	BuildVariant      string
	Command           string
	Count             string
	Configuration     string
	FileName          string
	FilePath          string
	Name              string
	Path              string
	Project           string
	SourceFile        string
	Target            string
	TargetFile        string
	TestCase          string
	TestFailureReason string
	TestSuite         string
	Time              string
}

func (e *entry) Match(txt string) (bool, map[string]string) {
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

func FillStruct(data map[string]string, result interface{}) {
	t := reflect.ValueOf(result).Elem()
	for k, v := range data {
		val := t.FieldByName(k)
		if val.CanSet() {
			val.Set(reflect.ValueOf(v))
		}
	}
}

func Parse(r io.Reader) {
	matchers := NewMatcher(simplereporter{})

	entry := LogEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()

		for _, matcher := range matchers {
			b, m := matcher.Match(txt)
			if b {
				mapstructure.Decode(m, &entry)
				matcher.logfunc(entry)
				break
			}
		}

	}
}
