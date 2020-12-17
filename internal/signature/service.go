package signature

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode/pbx"
	"fmt"

	"github.com/rs/zerolog/log"
)

const (
	KeyDevelopmentTeam  = "DEVELOPMENT_TEAM"
	KeyProfileSpecifier = "PROVISIONING_PROFILE_SPECIFIER"
	KeySigningIdentity  = "CODE_SIGN_IDENTITY"
	KeySigningStyle     = "CODE_SIGN_STYLE"
	ManualSigning       = "Manual"
)

type service struct {
	api.API
	*api.Config
}

func New(api api.API, cfg *api.Config) api.SignatureService {
	return service{api, cfg}
}

func (s service) Run(ctx context.Context) error {
	var res []api.TargetSignatureConfig

	// Parsing project
	pj, err := s.API.XCodeProjectService().Parse(ctx)
	if err != nil {
		return err
	}

	// Found configuration, installing it into a temporary keychain
	if err := s.API.KeyChainService().Create(ctx, "dothething"); err != nil {
		return err
	}

	if err = s.forTarget(ctx, s.Config.Target, pj, &res); err != nil {
		return err
	}

	for _, e := range res {
		s.applyTargetConfiguration(ctx, pj, e.TargetName, e.Config)
	}

	return err
}

// configureBuildSetting will apply the build settings for the XCBuildConfiguration
func (s service) configureBuildSetting(
	ctx context.Context,
	cc pbx.XCBuildConfiguration,
	m map[string]string,
) error {
	var err error
	for key, value := range m {
		path := fmt.Sprintf("buildSettings:%v", key)

		val, hasKey := cc.BuildSettings[key]
		if hasKey {
			err = s.API.PListBuddyService().SetStringValue(ctx, cc.Reference, path, value)
		} else {
			err = s.API.PListBuddyService().AddStringValue(ctx, cc.Reference, path, value)
		}
	}

	return err
}

func (a service) applyTargetConfiguration(
	ctx context.Context,
	pj api.Project,
	targetName string,
	sc *api.SignatureConfiguration,
) error {
	log.Debug().Str("Target", targetName).Msg("Configure target")

	// resolve target
	tgt, err := pj.Pbx.FindTargetByName(targetName)
	if err != nil {
		return NewSignatureError(err, ErrorTargetResolution)
	}

	// resolve configuration
	bc, err := tgt.BuildConfigurationList.FindConfiguration(a.Configuration)
	if err != nil {
		return NewSignatureError(err, ErrorBuildConfigurationResolution)
	}

	// installing provisioning profile
	err = a.API.ProvisioningService().Install(sc.ProvisioningProfile)
	if err != nil {
		return NewSignatureError(err, ErrorProvisioningInstall)
	}

	if err = a.configureBuildSetting(
		ctx,
		bc,
		map[string]string{
			KeyDevelopmentTeam:  sc.ProvisioningProfile.Entitlements.TeamID,
			KeyProfileSpecifier: sc.ProvisioningProfile.UUID,
			KeySigningIdentity:  sc.Cert.Issuer.CommonName,
			KeySigningStyle:     ManualSigning,
		}); err != nil {
		return err
	}

	if err = a.API.
		KeyChainService().
		ImportCertificate(ctx, sc.Cert.FilePath, a.CodeSignOption.CertificatePassword); err != nil {
		return NewSignatureError(err, ErrorCertificateImport)
	}

	return nil
}

func (s service) forTarget(
	ctx context.Context,
	t string,
	p api.Project,
	res *[]api.TargetSignatureConfig,
) error {
	// Resolving target by name
	nt, err := p.Pbx.FindTargetByName(t)
	if err != nil {
		return err
	}

	// Do the native target has native depdendencies
	for _, dp := range nt.Dependencies {
		switch dp.ProductType {
		case pbx.Application, pbx.TvExtension:
			if err := s.forTarget(ctx, dp.Name, p, res); err != nil {
				return err
			}
		}
	}

	// Resolving the build configuration for the target
	bc, err := nt.BuildConfigurationList.FindConfiguration(s.Config.Configuration)
	if err != nil {
		return err
	}

	// Resolving signature configuration for the bundle identifier
	sc, err := s.API.
		SignatureResolver().
		Resolve(ctx, bc.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"], nt.ProductType)
	if err != nil {
		return err
	}

	*res = append(*res, api.TargetSignatureConfig{
		TargetName: t,
		Config:     sc,
	})

	return nil
}
