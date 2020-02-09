package xcode

import (
	"dothething/internal/util"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	xCodeBuild = "xcodebuild"

	flagList = "-list"

	flagJSON = "-json"
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
}

// XCodeProjectService struct definition
type XCodeProjectService struct {
	exec util.Exec
}

// NewProjectService Create a new instance of the project service
func NewProjectService(exec util.Exec) *XCodeProjectService {
	return &XCodeProjectService{exec: exec}
}

// Parse the project
func (s *XCodeProjectService) Parse(projectPath *string) (*Project, error) {
	// Execute the list call to xcodebuild
	data, err := s.xCodeCall(projectPath)
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

func (s *XCodeProjectService) xCodeCall(path *string) ([]byte, error) {
	return s.exec.Exec(path, xCodeBuild, flagList, flagJSON)
}
