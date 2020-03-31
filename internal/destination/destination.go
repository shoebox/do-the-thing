package destination

import (
	"bufio"
	"context"
	"errors"

	"dothething/internal/util"
	"dothething/internal/xcode"
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
type DestinationService interface {
	Boot(ctx context.Context, d Destination) error
	List(scheme string) ([]Destination, error)
	ShutDown(ctx context.Context, d Destination) error
}

type destinationService struct {
	xcode xcode.XCodeBuildService
	exec  util.Exec
}

// NewDestinationService Create a new instance of the project service
func NewDestinationService(service xcode.XCodeBuildService, exec util.Exec) DestinationService {
	return destinationService{exec: exec, xcode: service}
}

// Boot boot a destination
func (s destinationService) Boot(ctx context.Context, d Destination) error {

	errc := make(chan error, 1)
	resc := make(chan string, 1)

	log.Info().
		Str("Destination ID", d.Id).
		Msg("Waiting for boot of destination")

	go s.xcRun(ctx, resc, errc, actionBootStatus, d.Id, flagBoot)

	select {
	case err := <-errc: // Checking for error
		return err

	case res := <-resc: // Resolving result
		log.Debug().Str("Result", res).Msg("Booting results")
		return nil

	case <-ctx.Done():
		if err := ctx.Err(); err != nil { // Checking for timeout
			return err
		}
	}

	return nil
}

func (s destinationService) xcRun(ctx context.Context,
	resc chan string,
	errc chan error,
	args ...string) {

	a := append([]string{simCtl}, args...)

	go func() {
		b, err := s.exec.ContextExec(ctx, xcRun, a...)
		if err != nil {
			errc <- err
			resc <- ""
		} else {
			resc <- string(b)
			errc <- nil
		}
	}()
}

// ShutDown a device
func (s destinationService) ShutDown(ctx context.Context, d Destination) error {
	log.Info().Str("Destination ID", d.Id).Msg("Shutdown destination")
	if _, err := s.exec.ContextExec(ctx, xcRun, simCtl, actionShutdown, d.Id); err != nil {
		return err
	}

	return nil
}

// ListDestinations Lists the valid destinations for a project or workspace and scheme
func (s destinationService) List(scheme string) ([]Destination, error) {
	res, err := s.xcode.ShowDestinations(scheme)
	if err != nil {
		return nil, ErrDestinationResolutionFailed
	}

	return s.parseDestinations(res), nil
}

func (s destinationService) parseDestinations(data string) []Destination {
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
