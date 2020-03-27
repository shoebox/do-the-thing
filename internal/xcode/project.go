package xcode

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Parse() (*Project, error)
}

// XCodeProjectService struct definition
type XCodeProjectService struct {
	xcodeService XCodeBuildService
}

// NewProjectService Create a new instance of the project service
func NewProjectService(service XCodeBuildService) ProjectService {
	return XCodeProjectService{xcodeService: service}
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
func (s XCodeProjectService) xCodeCall() ([]byte, error) {
	str, err := s.xcodeService.List()
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
