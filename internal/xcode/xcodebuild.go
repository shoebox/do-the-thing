package xcode

import (
	"context"
	"dothething/internal/util"
	"path/filepath"
)

const (
	// XCodeBuild executable
	Build = "xcodebuild"

	// FlagList list
	FlagList = "-list"

	FlagJSON = "-json"

	// FlagDestination destination specifier describing the device (or devices) to use as a destination
	FlagDestination = "-destination"

	// FlagShowDestinations Lists the valid destinations for a project or workspace and scheme.
	FlagShowDestinations = "-showdestinations"

	// FlagScheme Build the scheme specified by scheme name
	FlagScheme = "-scheme"

	// FlagProject Build the designated project
	FlagProject = "-project"

	// FlagWorkspace Build the designated workspace
	FlagWorkspace = "-workspace"

	FlagParallelTesting = "-parallel-testing-enabled"

	FlagParallelWorkerCount = "-maximum-parallel-testing-workers"

	// FlagResultBundlePath Writes a bundle to the specified path with results from performing an
	// action on a scheme in a workspace
	FlagResultBundlePath = "-resultBundlePath"

	// ActionTest Test a scheme from the build root
	ActionTest = "test"

	// ActionClean Remove build products and intermediate files from the build root
	ActionClean = "clean"
)

// XCodeBuildService service definition
type BuildService interface {
	List(ctx context.Context) (string, error)
	ShowDestinations(ctx context.Context, scheme string) (string, error)
	GetArg() string
	GetProjectPath() string
}

type xcodeBuildService struct {
	exec        util.Executor
	arg         string
	projectPath string
}

// NewService creates a new instance of the xcodebuild service
func NewService(exec util.Executor, projectPath string) BuildService {
	arg := FlagProject
	if filepath.Ext(projectPath) == ".xcworkspace" {
		arg = FlagWorkspace
	}
	return xcodeBuildService{exec: exec, arg: arg, projectPath: projectPath}
}

func (s xcodeBuildService) GetArg() string {
	return s.arg
}

func (s xcodeBuildService) GetProjectPath() string {
	return s.projectPath
}

// List Lists the targets and configurations in a project, or the schemes in a workspace
func (s xcodeBuildService) List(ctx context.Context) (string, error) {
	cmd := s.exec.CommandContext(ctx, Build, FlagList, FlagJSON, s.arg, s.projectPath)
	b, err := cmd.Output()
	if err != nil {
		return "", ParseXCodeBuildError(err)
	}

	return string(b), nil
}

func (s xcodeBuildService) ShowDestinations(ctx context.Context, scheme string) (string, error) {
	cmd := s.exec.CommandContext(ctx,
		Build,
		FlagShowDestinations,
		s.arg,
		s.projectPath,
		FlagScheme,
		scheme)

	b, err := cmd.Output()
	return string(b), ParseXCodeBuildError(err)
}
