package xcode

import (
	"context"
	"dothething/internal/api"
	"path/filepath"
)

const (
	ActionArchive           = "archive"                           // ActionArchive
	ActionBuild             = "build"                             // ActionBuild
	ActionClean             = "clean"                             // ActionClean Remove build products and intermediate files from the build root
	ActionTest              = "test"                              // ActionTest Test a scheme from the build root
	Cmd                     = "xcodebuild"                        // XCodeBuild executable
	FlagDestination         = "-destination"                      // FlagDestination destination specifier describing the device (or devices) to use as a destination
	FlagJSON                = "-json"                             // FlagJSON
	FlagList                = "-list"                             // FlagList list
	FlagParallelTesting     = "-parallel-testing-enabled"         // FlagParallelTesting
	FlagParallelWorkerCount = "-maximum-parallel-testing-workers" // FlagParallelWorkerCount
	FlagProject             = "-project"                          // FlagProject Build the designated project
	FlagResultBundlePath    = "-resultBundlePath"                 // FlagResultBundlePath Writes a bundle to the specified path with results from performing an action on a scheme in a workspace
	FlagScheme              = "-scheme"                           // FlagScheme Build the scheme specified by scheme name
	FlagShowDestinations    = "-showdestinations"                 // FlagShowDestinations Lists the valid destinations for a project or workspace and scheme.
	FlagConfiguration       = "-configuration"
	FlagWorkspace           = "-workspace" // FlagWorkspace Build the designated workspace
	FlagArchivePath         = "-archivePath"
	FlagExportPath          = "-exportPath"
)

type xcodeBuildService struct {
	*api.API
}

// NewService creates a new instance of the xcodebuild service
func NewService(api *api.API) api.BuildService {
	return xcodeBuildService{API: api}
}

// GetArg returns the right flag to execute depending of the type of project path configured
func (s xcodeBuildService) GetArg() string {
	arg := FlagProject
	if filepath.Ext(s.API.Config.Path) == ".xcworkspace" {
		arg = FlagWorkspace
	}
	return arg
}

// List lists the targets and configurations in a project, or the schemes in a workspace
func (s xcodeBuildService) List(ctx context.Context) (string, error) {
	// Executing command
	cmd := s.API.Exec.CommandContext(ctx, Cmd, FlagList, FlagJSON, s.GetArg(), s.API.Config.Path)

	// Resolving combined outputs
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", ParseXCodeBuildError(err)
	}

	return string(b), nil
}

// ShowDestinations will resolve the destinations for the scheme
func (s xcodeBuildService) ShowDestinations(ctx context.Context, scheme string) (string, error) {
	cmd := s.API.Exec.CommandContext(ctx,
		Cmd,
		FlagShowDestinations,
		s.GetArg(),
		s.API.Config.Path,
		FlagScheme,
		scheme)

	b, err := cmd.Output()

	return string(b), ParseXCodeBuildError(err)
}
