package project

import (
	"bytes"
	"context"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/pbx"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
)

var (
	// ErrInvalidConfig The xcodebuild answer was not valid
	ErrInvalidConfig = errors.New("Invalid xcodedbuild list response")
)

// Project datas
type Project struct {
	Configurations []string `json:"configurations"`
	Name           string   `json:"name"`
	Pbx            pbx.PBXProject
	Schemes        []string `json:"schemes"`
	Targets        []string `json:"targets"`
}

// ProjectService interface
type ProjectService interface {
	Parse(ctx context.Context) (Project, error)
}

// projectService struct definition
type projectService struct {
	xcodeService xcode.BuildService
	exec         util.Executor
}

type list struct {
	Project   Project `json:"project"`
	Workspace Project `json:"workspace"`
}

// NewProjectService Create a new instance of the PBXProj service
func NewProjectService(service xcode.BuildService, e util.Executor) ProjectService {
	return projectService{service, e}
}

// Parse the PBXProj raw project into a more friendly version
func (s projectService) Parse(ctx context.Context) (Project, error) {
	project, err := s.resolveProject(ctx)
	if err != nil {
		return Project{}, err
	}

	project.Pbx, err = s.resolvePbx()
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

func (s projectService) resolvePbx() (pbx.PBXProject, error) {
	var res pbx.PBXProject

	// Resolve the project path
	b, err := s.resolvePbxProj()
	if err != nil {
		return res, err
	}

	return s.decodeProject(b)
}

func (s projectService) decodeProject(b []byte) (pbx.PBXProject, error) {
	var raw pbx.PBXProjRaw
	if err := util.DecodeFile(bytes.NewReader(b), &raw); err != nil {
		return pbx.PBXProject{}, err
	}
	return raw.Parse(), nil
}

func (s projectService) resolvePbxProj() ([]byte, error) {
	path, err := filepath.Abs(s.xcodeService.GetProjectPath() + "/project.pbxproj")
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(path)
}

func (s projectService) resolveProject(ctx context.Context) (Project, error) {
	var project Project

	// Execute the list call to xcodebuild
	data, err := s.xCodeCall(ctx)
	if err != nil {
		return project, err
	}

	// Unmarshall the response
	var root list
	err = json.Unmarshal(data, &root)
	if err != nil {
		return project, ErrInvalidConfig
	}

	if root.Project.Name != "" {
		return root.Project, nil
	}

	return root.Workspace, nil
}

func (s projectService) xCodeCall(ctx context.Context) ([]byte, error) {
	str, err := s.xcodeService.List(ctx)
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
