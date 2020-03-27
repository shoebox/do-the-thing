package destination

import (
	"bufio"
	"errors"

	"dothething/internal/xcode"
	"reflect"
	"regexp"
	"strings"
)

// ErrDestinationResolutionFailed Failed to resolve destinations for the project
var ErrDestinationResolutionFailed = errors.New("Command execution failed")

// Destination available destination for the scheme
type Destination struct {
	Name     string
	Platform string
	Id       string
	OS       string
}

type DestinationService interface {
	Boot(d Destination)
	List(scheme string) ([]Destination, error)
}

type XCodeDestinationService struct {
	xcode xcode.XCodeBuildService
}

// NewDestinationService Create a new instance of the project service
func NewDestinationService(service xcode.XCodeBuildService) DestinationService {
	return XCodeDestinationService{xcode: service}
}

// Boot boot a destination
func (s XCodeDestinationService) Boot(d Destination) {
}

// InstallOnDestination Install an app on a device
func (s XCodeDestinationService) Install(d Destination, path string) {
}

// LaunchOnDestination launch an application by identifier on a device
func (s XCodeDestinationService) Launch(d Destination, id string) {
}

// ShutDown a device
func (s XCodeDestinationService) ShutDown(d Destination) {
}

// ListDestinations Lists the valid destinations for a project or workspace and scheme
func (s XCodeDestinationService) List(scheme string) ([]Destination, error) {
	res, err := s.xcode.ShowDestinations(scheme)
	if err != nil {
		return nil, ErrDestinationResolutionFailed
	}

	return s.parseDestinations(res), nil
}

func (s XCodeDestinationService) parseDestinations(data string) []Destination {
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
