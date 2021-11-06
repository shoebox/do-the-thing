package project

import (
	"bytes"
	"context"
	"dothething/internal/api"
	"dothething/internal/util"
	"dothething/internal/xcode/pbx"
	"encoding/json"
	"errors"
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

func (p Project) ValidateConfiguration(c api.Config) error {
	found := false
	for _, s := range p.Schemes {
		if s == c.Scheme {
			found = true
			break
		}
	}

	if !found {
		return errors.New("Invalid scheme for the project")
	}

	return nil
}

// projectService struct definition
type projectService struct {
	*api.API
}

type list struct {
	Project   api.Project `json:"project"`
	Workspace api.Project `json:"workspace"`
}

// NewProjectService Create a new instance of the PBXProj service
func NewProjectService(api *api.API) api.ProjectService {
	return projectService{api}
}

// Parse the PBXProj raw project into a more friendly version
func (s projectService) Parse(ctx context.Context) (api.Project, error) {
	var res api.Project
	res, err := s.resolveProject(ctx)
	if err != nil {
		return res, err
	}

	res.Pbx, err = s.resolvePbx()
	if err != nil {
		return res, err
	}

	return res, nil
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
	path, err := filepath.Abs(s.API.PathService.XCodeProject() + "/project.pbxproj")
	if err != nil {
		return []byte{}, err
	}

	return s.API.FileService.OpenAndReadFileContent(path)
}

func (s projectService) resolveProject(ctx context.Context) (api.Project, error) {
	var project api.Project

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
	str, err := s.API.BuildService.List(ctx)
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
