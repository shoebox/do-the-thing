package action

import (
	"context"
	"dothething/internal/config"
	"dothething/internal/util"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

type ActionArchive interface {
	Run(ctx context.Context, config config.Config) error
}

func NewArchive(xcode xcode.BuildService, exec util.Executor) ActionArchive {
	return actionArchive{exec, xcode}
}

type actionArchive struct {
	exec  util.Executor
	xcode xcode.BuildService
}

func (a actionArchive) Run(ctx context.Context, config config.Config) error {
	xce := xcode.ParseXCodeBuildError(a.build(ctx, config))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionArchive) build(ctx context.Context, config config.Config) error {
	log.Info().Msg("Archiving")
	return RunCmd(a.exec.CommandContext(ctx,
		xcode.Cmd,
		a.xcode.GetArg(),
		a.xcode.GetProjectPath(),
		xcode.ActionArchive,
		xcode.FlagScheme, config.Scheme,
		"CODE_SIGN_IDENTITY=iPhone developer: Self signer",
		"CODE_SIGN_STYLE=Manual",
		"DEVELOPMENT_TEAM=12345ABCDE",
		"PROVISIONING_PROFILE_SPECIFIER=B5C2906D-D6EE-476E-AF17-D99AE14644AA",

		// "-xcconfig", "/var/folders/m1/s05mrl8s4fbfw8zf4hw04tjwh740pl/T/do-the-thing-736615920/file.xcconfig",
	),
	)
	//"CODE_SIGNING_ALLOWED=NO"))
}
