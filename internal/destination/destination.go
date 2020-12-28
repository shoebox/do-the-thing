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
	cmd := s.API.Exec().CommandContext(ctx, xcRun, simCtl, actionBootStatus, d.ID, flagBoot)

	b, err := cmd.Output()
	if err != nil {
		return NewBootError(d.ID)
	}

	log.Info().
		Str("Result", string(b)).
		Msg("Booting results")

	// TODO: Handle booting results better

	return nil
}

// ShutDown a device
func (s destinationService) ShutDown(ctx context.Context, d api.Destination) error {
	log.Info().Str("Destination ID", d.ID).Msg("Shutdown destination")
	cmd := s.API.
		Exec().
		CommandContext(ctx, xcRun, simCtl, actionShutdown, d.ID)

	if _, err := cmd.Output(); err != nil {
		return NewShutDownError(d.ID)
	}

	return nil
}

// ListDestinations Lists the valid destinations for a project or workspace and scheme
func (s destinationService) List(ctx context.Context, scheme string) ([]api.Destination, error) {
	res, err := s.API.XCodeBuildService().ShowDestinations(ctx, scheme)
	if err != nil {
		return nil, NewListingError()
	}

	return s.parseDestinations(res), err
}

var destinationRegexp = regexp.MustCompile(`([^:\s]+):([^,}]+)[,]?`)

func (s destinationService) parseDestinations(data string) []api.Destination {
	// Result
	var res []api.Destination

	sc := bufio.NewScanner(strings.NewReader(data))
	start := false

	// For each output lines
	for sc.Scan() {
		s.parseLine(sc.Text(), &start, &res)
	}

	return res
}

func (s destinationService) parseLine(line string, start *bool, res *[]api.Destination) {
	if strings.Contains(line, "Available destinations") { // Start of section containing the available destinations
		*start = true
	} else if line == "" { // End of section
		*start = false
	} else if *start {
		// Split on regex
		indexes := destinationRegexp.FindAllSubmatch([]byte(line), -1)

		// Map the splitted values
		m := map[string]string{}
		for _, r := range indexes {
			name := string(r[1])

			// Special case, for the ID
			if name == "id" {
				name = "ID"
			}

			//
			m[name] = strings.TrimSpace(string(r[2]))
		}

		// Populate the destination
		dest := api.Destination{}
		fillStruct(m, &dest)

		// Append destination
		*res = append(*res, dest)
	}
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
