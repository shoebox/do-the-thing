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
	api.API
	cfg *api.Config
}

// NewService creates a new instance of the xcodebuild service
func NewService(api api.API, cfg *api.Config) api.BuildService {
	return xcodeBuildService{API: api, cfg: cfg} //, arg: arg, projectPath: projectPath}
}

func (s xcodeBuildService) GetArg() string {
	arg := FlagProject
	if filepath.Ext(s.cfg.Path) == ".xcworkspace" {
		arg = FlagWorkspace
	}
	return arg
}

func (s xcodeBuildService) GetProjectPath() string {
	return s.cfg.Path
}

// List Lists the targets and configurations in a project, or the schemes in a workspace
func (s xcodeBuildService) List(ctx context.Context) (string, error) {
	cmd := s.API.Exec().CommandContext(ctx, Cmd, FlagList, FlagJSON, s.GetArg(), s.cfg.Path)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", ParseXCodeBuildError(err)
	}

	return string(b), nil
}

func (s xcodeBuildService) ShowDestinations(ctx context.Context, scheme string) (string, error) {
	cmd := s.API.Exec().CommandContext(ctx,
		Cmd,
		FlagShowDestinations,
		s.GetArg(),
		s.cfg.Path,
		FlagScheme,
		scheme)

	b, err := cmd.Output()
	return string(b), ParseXCodeBuildError(err)
}
