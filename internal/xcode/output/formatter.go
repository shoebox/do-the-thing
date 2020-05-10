package output

import (
	"bufio"
	"io"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

var matchers = NewMatcher(simplereporter{})

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
