package action

import (
	"context"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/output"

	"github.com/rs/zerolog/log"
)

type ActionBuild interface {
	Run(ctx context.Context, config xcode.Config) error
}

func NewBuild(xcode xcode.XCodeBuildService, exec util.Executor) ActionBuild {
	return actionBuild{exec, xcode}
}

type actionBuild struct {
	exec  util.Executor
	xcode xcode.XCodeBuildService
}

func (a actionBuild) Run(ctx context.Context, config xcode.Config) error {

	log.Info().
		Msg("ActionArchive")

	cmd := a.exec.CommandContext(ctx,
		xcode.XCodeBuild,
		a.xcode.GetArg(),
		a.xcode.GetProjectPath(),
		xcode.ActionClean,
		"build",
		xcode.FlagScheme, config.Scheme,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO")

	pout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	perr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		f := output.NewFormatter(output.SimpleReporter{})
		f.Parse(pout)
		f.Parse(perr)
	}()

	//
	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
