package destination

import (
	"bufio"
	"context"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	actionBootStatus = "bootstatus"
	actionShutdown   = "shutdown"
	simCtl           = "simctl"
	xcRun            = "xcrun"

	flagBoot = "-b" // Boot the device if it isn't already booted
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

// DestinationService destination service definition
type Service interface {
	Boot(ctx context.Context, d Destination) error
	List(ctx context.Context, scheme string) ([]Destination, error)
	ShutDown(ctx context.Context, d Destination) error
}

type destinationService struct {
	xcode xcode.BuildService
	exec  util.Executor
}

// NewDestinationService Create a new instance of the project service
func NewDestinationService(service xcode.BuildService, exec util.Executor) Service {
	return destinationService{exec: exec, xcode: service}
}

// Boot boot a destination
func (s destinationService) Boot(ctx context.Context, d Destination) error {
	cmd := s.exec.CommandContext(ctx, xcRun, simCtl, actionBootStatus, d.Id, flagBoot)

	b, err := cmd.Output()
	if err != nil {
		return err
	}

	log.Info().
		Str("Result", string(b)).
		Msg("Booting results")

	return nil
}

// ShutDown a device
func (s destinationService) ShutDown(ctx context.Context, d Destination) error {
	log.Info().Str("Destination ID", d.Id).Msg("Shutdown destination")
	cmd := s.exec.CommandContext(ctx, xcRun, simCtl, actionShutdown, d.Id)

	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}

// ListDestinations Lists the valid destinations for a project or workspace and scheme
func (s destinationService) List(ctx context.Context, scheme string) ([]Destination, error) {
	res, err := s.xcode.ShowDestinations(ctx, scheme)
	return s.parseDestinations(res), err
}

func (s destinationService) parseDestinations(data string) []Destination {
	// Result
	var res []Destination

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
