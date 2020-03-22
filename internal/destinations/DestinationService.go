package destinations

import (
	"context"
	"dothething/internal/util"
)

type Service interface {
	List() ([]Destination, error)
	Boot(d Destination)
}

type Destination struct {
	Platform string
	OS       string
	Name     string
	ID       string
}

type DestinationService struct {
	exec util.Exec
}

func (s DestinationService) List(projectPath string, scheme string) ([]Destination, error) {
	s.exec.Exec(nil, "xcodebuild", "-showdestinations", "-project", projectPath, "-scheme", scheme)

	return []Destination{}, nil
}

func (s Destination) Boot(d Destination) {

}

func (s DestinationService) waitForDestinationStatus(c context.Context, d Destination, status string) error {
	return nil
}
