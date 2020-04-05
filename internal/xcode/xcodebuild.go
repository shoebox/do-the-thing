package xcode

import (
	"context"
	"dothething/internal/util"
	"path/filepath"
)

const (
	// XCodeBuild executable
	XCodeBuild = "xcodebuild"

	// FlagList list
	flagList = "-list"

	flagJSON = "-json"

	// FlagShowDestinations Lists the valid destinations for a project or workspace and scheme.
	FlagShowDestinations = "-showdestinations"

	// FlagScheme Build the scheme specified by scheme name
	FlagScheme = "-scheme"

	// FlagProject Build the designated project
	FlagProject = "-project"

	// FlagWorkspace Build the designated workspace
	FlagWorkspace = "-workspace"
)

// XCodeBuildService service definition
type XCodeBuildService interface {
	List(ctx context.Context) (string, error)
	ShowDestinations(ctx context.Context, scheme string) (string, error)
	Run(ctx context.Context, arg ...string) (string, error)
}

type xcodeBuildService struct {
	exec        util.Exec
	arg         string
	projectPath string
}

// NewService creates a new instance of the xcodebuild service
func NewService(exec util.Exec, projectPath string) XCodeBuildService {
	arg := FlagProject
	if filepath.Ext(projectPath) == ".xcworkspace" {
		arg = FlagWorkspace
	}
	return xcodeBuildService{exec: exec, arg: arg, projectPath: projectPath}
}

// List Lists the targets and configurations in a project, or the schemes in a workspace
func (s xcodeBuildService) List(ctx context.Context) (string, error) {
	return s.Run(ctx, flagList, flagJSON, s.arg, s.projectPath)
}

func (s xcodeBuildService) ShowDestinations(ctx context.Context, scheme string) (string, error) {
	return s.Run(ctx, FlagShowDestinations, s.arg, s.projectPath, FlagScheme, scheme)
}

func (s xcodeBuildService) Run(ctx context.Context, arg ...string) (string, error) {
	errc := make(chan error, 1)
	resc := make(chan string, 1)

	// Execute command
	go func() {
		b, err := s.exec.ContextExec(ctx,
			XCodeBuild,
			arg...)
		if err != nil {
			errc <- err
		} else {
			resc <- string(b)
		}
	}()

	select {
	case err := <-errc: // Checking for error
		return "", err

	case res := <-resc: // Resolving result
		return res, nil

	case <-ctx.Done():
		if err := ctx.Err(); err != nil { // Checking for timeout
			return "", err
		}
	}

	return "", nil
}
