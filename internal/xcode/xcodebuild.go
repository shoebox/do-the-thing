package xcode

import (
	"context"
	"dothething/internal/util"
	"path/filepath"
	"time"
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

type XCodeBuildService interface {
	List() (string, error)
	ShowDestinations(scheme string) (string, error)
	Run(arg ...string) (string, error)
}

type XCodeBuildServiceImpl struct {
	exec        util.Exec
	arg         string
	projectPath string
}

func NewService(exec util.Exec, projectPath string) XCodeBuildService {
	arg := FlagProject
	if filepath.Ext(projectPath) == ".xcworkspace" {
		arg = FlagWorkspace
	}
	return XCodeBuildServiceImpl{exec: exec, arg: arg, projectPath: projectPath}
}

// List Lists the targets and configurations in a project, or the schemes in a workspace
func (s XCodeBuildServiceImpl) List() (string, error) {
	return s.Run(flagList, flagJSON, s.arg, s.projectPath)
}

func (s XCodeBuildServiceImpl) ShowDestinations(scheme string) (string, error) {
	return s.Run(FlagShowDestinations, s.arg, s.projectPath, FlagScheme, scheme)
}

func (s XCodeBuildServiceImpl) Run(arg ...string) (string, error) {
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

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
