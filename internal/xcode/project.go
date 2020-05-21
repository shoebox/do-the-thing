package xcode

import (
	"context"
	"encoding/json"
	"errors"
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

// ProjectService interface
type ProjectService interface {
	Parse(ctx context.Context) (*Project, error)
}

// projectService struct definition
type projectService struct {
	xcodeService BuildService
}

// NewProjectService Create a new instance of the project service
func NewProjectService(service BuildService) ProjectService {
	return projectService{xcodeService: service}
}

// Parse the project
func (s projectService) Parse(ctx context.Context) (*Project, error) {
	// Execute the list call to xcodebuild
	data, err := s.xCodeCall(ctx)
	if err != nil {
		return nil, err
	}

	// Unmarshall the response
	var root list
	err = json.Unmarshal(data, &root)
	if err != nil {
		return nil, ErrInvalidConfig
	}

	return &root.Project, nil
}

func (s projectService) xCodeCall(ctx context.Context) ([]byte, error) {
	str, err := s.xcodeService.List(ctx)
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
