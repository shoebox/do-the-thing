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

	flagList = "-list"

	// Lists the valid destinations for a project or workspace and scheme.
	flagShowDestinations = "-showdestinations"

	// Build the scheme specified by scheme name
	flagScheme = "-scheme"

	// Json Output
	flagJSON = "-json"

	// Build the project
	flagProject = "-project"

	// Build the workspace
	flagWorkspace = "-workspace"
)

var (
	// ErrInvalidConfig The xcodebuild answer was not valid
	ErrInvalidConfig = errors.New("Invalid xcodedbuild list response")

	// ErrDestinationResolutionFailed Failed to resolve destinations for the project
	ErrDestinationResolutionFailed = errors.New("Command execution failed")
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
	if filepath.Ext(projectPath) == ".xcworkspace" {
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

// ListDestinations Lists the valid destinations for a project or workspace and scheme
func (s XCodeProjectService) ListDestinations(scheme string) ([]Destination, error) {

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	errc := make(chan error, 1)
	resc := make(chan string, 1)

	// Execute command
	go func() {
		b, err := s.exec.ContextExec(ctx,
			xCodeBuild,
			flagShowDestinations,
			s.arg, s.path,
			flagScheme, scheme)
		if err != nil {
			errc <- err
		} else {
			resc <- string(b)
		}
	}()

	select {
	case err := <-errc: // Checking for error
		return nil, err

	case res := <-resc: // Resolving result
		return s.parseDestinations(res), nil

	case <-ctx.Done():
		if err := ctx.Err(); err != nil { // Checking for timeout
			return nil, ErrDestinationResolutionFailed
		}
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
				m[string(r[1])] = strings.TrimSpace(string(r[2]))
			}

			// Populate the destination
			dest := Destination{}
			fillStruct(m, &dest)

			// Append destination
			res = append(res, dest)
		}
	}

	return res
}

func fillStruct(data map[string]string, result interface{}) interface{} {
	t := reflect.ValueOf(result)
	elem := t.Elem()

	for k, v := range data {
		name := strings.Title(k)
		field := elem.FieldByName(name)
		if field.IsValid() {
			val := elem.FieldByName(name)
			val.Set(reflect.ValueOf(v))
		}
	}

	return result
}

func (s XCodeProjectService) xCodeCall() ([]byte, error) {
	return s.exec.Exec(nil, xCodeBuild, flagList, flagJSON, s.arg, s.path)
}
