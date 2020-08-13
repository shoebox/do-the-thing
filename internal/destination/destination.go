package destination

import (
	"bufio"
	"context"
	"dothething/internal/api"
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

type destinationService struct {
	api.API
}

// NewDestinationService Create a new instance of the project service
func NewDestinationService(a api.API) api.DestinationService {
	return destinationService{a}
}

// Boot boot a destination
func (s destinationService) Boot(ctx context.Context, d api.Destination) error {
	cmd := s.API.Exec().CommandContext(ctx, xcRun, simCtl, actionBootStatus, d.Id, flagBoot)

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
func (s destinationService) ShutDown(ctx context.Context, d api.Destination) error {
	log.Info().Str("Destination ID", d.Id).Msg("Shutdown destination")
	cmd := s.API.Exec().CommandContext(ctx, xcRun, simCtl, actionShutdown, d.Id)

	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}

// ListDestinations Lists the valid destinations for a project or workspace and scheme
func (s destinationService) List(ctx context.Context, scheme string) ([]api.Destination, error) {
	res, err := s.API.XCodeBuildService().ShowDestinations(ctx, scheme)
	return s.parseDestinations(res), err
}

func (s destinationService) parseDestinations(data string) []api.Destination {
	// Result
	var res []api.Destination

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
			dest := api.Destination{}
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
