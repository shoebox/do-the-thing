package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"
	"dothething/internal/xcode/pbx"
	"fmt"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

func NewArchive(api api.API, cfg *api.Config) api.Action {
	return actionArchive{api, cfg}
}

type actionArchive struct {
	api.API
	*api.Config
}

func (a actionArchive) Run(ctx context.Context) error {
	xce := xcode.ParseXCodeBuildError(a.build(ctx))
	fmt.Println("run", xce)
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionArchive) configureBuildSetting(
	ctx context.Context,
	cc pbx.XCBuildConfiguration,
	key string,
	value string,
) error {
	var err error
	bs := fmt.Sprintf("buildSettings:%v", key)

	_, ok := cc.BuildSettings[key]
	if ok {
		err = a.API.PListBuddyService().SetStringValue(ctx, cc.Reference, bs, value)
	} else {
		err = a.API.PListBuddyService().AddStringValue(ctx, cc.Reference, bs, value)
	}

	return err
}

func (a actionArchive) build(ctx context.Context) error {
	log.Info().Msg("Archiving")

	// Parsing project
	pj, err := a.API.XCodeProjectService().Parse(ctx)
	if err != nil {
		return err
	}

	// Resolving signature
	cfg, err := a.API.SignatureService().Run(ctx, pj)
	if err != nil {
		return err
	}

	for _, e := range cfg {
		tgt, err := pj.Pbx.FindTargetByName(e.TargetName)
		if err != nil {
			log.Error().Err(err).Msg("Error while resolving target")
			continue
		}

		cc, err := tgt.BuildConfigurationList.FindConfiguration(a.Configuration)
		if err != nil {
			log.Error().Err(err).Str("Build configuration", a.Configuration).Msg("Failed to resolve build configuration")
		}

		// Instal
		err = a.API.ProvisioningService().Install(e.Config.ProvisioningProfile)
		if err != nil {
			log.Error().Err(err).Msg("Error while installing certificate")
		}

		a.configureBuildSetting(ctx, cc, "DEVELOPMENT_TEAM", e.Config.ProvisioningProfile.Entitlements.TeamID)
		a.configureBuildSetting(ctx, cc, "PROVISIONING_PROFILE_SPECIFIER", e.Config.ProvisioningProfile.UUID)
		a.configureBuildSetting(ctx, cc, "CODE_SIGN_IDENTITY", e.Config.Cert.Issuer.CommonName)
		a.configureBuildSetting(ctx, cc, "CODE_SIGN_STYLE", "Manual")

		err = a.API.KeyChainService().
			ImportCertificate(ctx, e.Config.Cert.FilePath, a.CodeSignOption.CertificatePassword)
		fmt.Println("importing cert", err)
	}

	// The archiving arguments
	args := []string{
		a.API.XCodeBuildService().GetArg(),
		a.API.XCodeBuildService().GetProjectPath(),
		xcode.ActionArchive,
		xcode.FlagScheme, a.Config.Scheme,
		xcode.FlagConfiguration, a.Config.Configuration,
		xcode.FlagArchivePath, fmt.Sprintf("%v/archive/toto.xcarchive", a.Config.Path),
	}

	// fmt.Println(pj.Name, cfg)
	defer a.API.KeyChainService().Delete(ctx)

	return RunCmd(a.API.Exec().CommandContext(ctx, xcode.Cmd, args...))
}
