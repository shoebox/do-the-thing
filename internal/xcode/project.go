package xcode

import (
	"bufio"
	"context"
	"dothething/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	xCodeBuild = "xcodebuild"

	flagList             = "-list"
	flagShowDestinations = "-showdestinations"
	flagScheme           = "-scheme"

	flagJSON      = "-json"
	flagProject   = "-project"
	flagWorkspace = "-workspace"
)

var (
	// ErrInvalidConfig The xcodebuild answer was not valid
	ErrInvalidConfig = errors.New("Invalid xcodedbuild list response")
)

type list struct {
	Project Project `json:"project"`
}

// Project datas
type Project struct {
	Configurations []string `json:"configurations"`
	Name           string   `json:"name"`
	Schemes        []string `json:"schemes"`
	Targets        []string `json:"targets"`
}

type Destination struct {
	Name     string
	Platform string
	Id       string
	OS       string
}

// ProjectService interface
type ProjectService interface {
	Parse() (*Project, error)
	ListDestinations(scheme string) ([]Destination, error)
}

// XCodeProjectService struct definition
type XCodeProjectService struct {
	arg  string
	exec util.Exec
	path string
}

// NewProjectService Create a new instance of the project service
func NewProjectService(exec util.Exec, projectPath string) ProjectService {
	arg := flagProject
	if filepath.Ext(projectPath) == "xcworkspace" {
		arg = flagWorkspace
	}
	return XCodeProjectService{exec: exec, path: projectPath, arg: arg}
}

// Parse the project
func (s XCodeProjectService) Parse() (*Project, error) {
	// Execute the list call to xcodebuild
	data, err := s.xCodeCall()
	if err != nil {
		return nil, fmt.Errorf("Failed to call xcode API (Error : %s)", err)
	}

	// Unmarshall the response
	var root list
	err = json.Unmarshal(data, &root)
	if err != nil {
		return nil, ErrInvalidConfig
	}

	return &root.Project, nil
}

// ListDestinations Lists the valid destinations for a project or workspace and sche
func (s XCodeProjectService) ListDestinations(scheme string) ([]Destination, error) {

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	errc := make(chan error, 1)
	resc := make(chan string, 1)

	// Execute command
	go func() {
		b, err := s.exec.ContextExec(ctx,
			xCodeBuild,
			flagShowDestinations,
			s.arg,
			s.path,
			flagScheme,
			scheme)
		resc <- string(b)
		errc <- err
	}()

	select {
	case res := <-resc:
		d := s.parseDestinations(res)
		for _, r := range d {
			fmt.Printf("%#v\n", r)
		}
		return s.parseDestinations(res), nil

	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			return nil, err
		}

	case err := <-errc:
		return nil, err

	}

	return nil, nil
}

func (s XCodeProjectService) parseDestinations(data string) []Destination {
	// Result
	res := []Destination{}

	sc := bufio.NewScanner(strings.NewReader(data))
	start := false
	regex := regexp.MustCompile(`([^:\s]+):([^,}]+)[,]?`)

	// For each output lines
	for sc.Scan() {
		if strings.Contains(sc.Text(), "Available destinations") { // Start of section containing the available destinations
			start = true
		} else if sc.Text() == "" { // End of section
			start = false
		} else if start {
			// Split on regex
			indexes := regex.FindAllSubmatch([]byte(sc.Text()), -1)

			// Map the splitted values
			m := map[string]string{}
			for _, r := range indexes {
				m[string(r[1])] = string(r[2])
			}

			// Populate the destination
			dest := Destination{}
			fillStruct(m, &dest)
			res = append(res, dest)
		}
	}

	return res
}

func fillStruct(data map[string]string, result interface{}) interface{} {
	t := reflect.ValueOf(result).Elem()
	for k, v := range data {
		val := t.FieldByName(strings.Title(k))
		val.Set(reflect.ValueOf(v))
	}

	return result
}

func (s XCodeProjectService) xCodeCall() ([]byte, error) {
	return s.exec.Exec(nil, xCodeBuild, flagList, flagJSON, s.arg, s.path)
}
