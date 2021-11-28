package signature

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode/pbx"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const (
	KeyDevelopmentTeam  = "DEVELOPMENT_TEAM"
	KeyProfileSpecifier = "PROVISIONING_PROFILE_SPECIFIER"
	KeySigningIdentity  = "CODE_SIGN_IDENTITY"
	KeySigningStyle     = "CODE_SIGN_STYLE"
	ManualSigning       = "Manual"
)

func NewSignatureService(api *api.API) signatureService {
	return signatureService{API: api}
}

type signatureService struct {
	*api.API
}

var cfg []api.TargetSignatureConfig

func (s signatureService) GetConfiguration() *[]api.TargetSignatureConfig {
	return &cfg
}

func (s signatureService) Run(ctx context.Context) error {
	// Parsing project
	pj, err := s.API.XCodeProjectService.Parse(ctx)
	if err != nil {
		return err
	}

	// Found configuration, installing it into a temporary keychain
	log.Info().Msg("Found configuration")
	err = s.API.KeyChain.Create(ctx, "dothething")
	if err != nil {
		log.Error().AnErr("Error", err).Msg("Failed to create keychain")
		return err
	}

	// Resolving for target
	if err = s.forTarget(ctx, s.API.Config.Target, pj); err != nil {
		return err
	}

	for _, e := range cfg {
		log.Info().Str("Name", e.TargetName).Msg("Configuring target")
		if err := s.applyTargetConfiguration(ctx, pj, e.TargetName, e.Config); err != nil {
			log.Error().
				AnErr("Error", err).
				Str("Target", e.TargetName).
				Msg("Failed to configure target")
			return err
		}
	}

	return err
}

// configureBuildSetting will apply the build settings for the XCBuildConfiguration
func (s signatureService) configureBuildSetting(
	ctx context.Context,
	cc pbx.XCBuildConfiguration,
	m map[string]string,
) error {
	var err error
	for key, value := range m {
		path := s.buildSettingsPath(key)

		_, hasKey := cc.BuildSettings[key]
		if hasKey {
			err = s.API.PlistBuddyService.SetStringValue(ctx, cc.Reference, path, value)
		} else {
			err = s.API.PlistBuddyService.AddStringValue(ctx, cc.Reference, path, value)
		}
	}

	return err
}

func (s signatureService) buildSettingsPath(key string) string {
	return fmt.Sprintf("buildSettings:%v", key)
}

func (a signatureService) applyTargetConfiguration(
	ctx context.Context,
	pj api.Project,
	targetName string,
	sc *api.SignatureConfiguration,
) error {
	log.Info().Str("Target", targetName).Msg("Configuring target")

	// resolve target
	tgt, err := pj.Pbx.FindTargetByName(targetName)
	if err != nil {
		return NewSignatureError(err, ErrorTargetResolution)
	}

	// resolve configuration
	bc, err := tgt.BuildConfigurationList.FindConfiguration(a.API.Config.Configuration)
	if err != nil {
		return NewSignatureError(err, ErrorBuildConfigurationResolution)
	}

	// installing provisioning profile
	err = a.API.ProvisioningService.Install(sc.ProvisioningProfile)
	if err != nil {
		return NewSignatureError(err, ErrorProvisioningInstall)
	}

	if err = a.configureBuildSettingsOfBuildConfiguration(
		ctx,
		bc,
		sc.ProvisioningProfile.Entitlements.TeamID,
		sc.ProvisioningProfile.UUID,
		sc.Cert.Issuer.CommonName,
	); err != nil {
		return NewSignatureError(err, ErrorBuildSettingsConfiguration)
	}

	path, err := filepath.Abs(sc.Cert.FilePath)
	if err != nil {
		return NewSignatureError(err, ErrorCertificateImport)
	}

	if err = a.API.KeyChain.ImportCertificate(
		ctx,
		path,
		a.API.Config.CodeSignOption.CertificatePassword,
		sc.Cert.Issuer.CommonName,
	); err != nil {
		return NewSignatureError(err, ErrorCertificateImport)
	}

	return nil
}

func (a signatureService) configureBuildSettingsOfBuildConfiguration(
	ctx context.Context,
	bc pbx.XCBuildConfiguration,
	teamID string,
	UUID string,
	identity string,
) error {
	log.Info().
		Str("TeamID", teamID).
		Str("UUD", UUID).
		Str("Identity", identity).
		Msg("Configuring build settingd")

	return a.configureBuildSetting(
		ctx,
		bc,
		map[string]string{
			KeyDevelopmentTeam:  teamID,
			KeyProfileSpecifier: UUID,
			KeySigningIdentity:  identity,
			KeySigningStyle:     ManualSigning,
		},
	)
}

func (s signatureService) forTarget(
	ctx context.Context,
	t string,
	p api.Project,
) error {
	log.Info().Str("Target", t).Msg("Resolving for target")
	// Resolving target by name
	nt, err := p.Pbx.FindTargetByName(t)
	if err != nil {
		return fmt.Errorf("failed to find target %v (%v)", t, err)
	}

	if err = s.configureDependencies(nt, func(name string) error {
		return s.forTarget(ctx, name, p)
	}); err != nil {
		return err
	}

	// Resolving the build configuration for the target
	bc, err := nt.BuildConfigurationList.FindConfiguration(s.API.Config.Configuration)
	if err != nil {
		return fmt.Errorf("failed to find build configuration %v (%v)", s.API.Config.Configuration, err)
	}

	// Resolving signature configuration for the bundle identifier
	bundleID := bc.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"]
	sc, err := s.API.
		SignatureResolver.
		Resolve(ctx, bundleID, nt.ProductType)
	if err != nil {
		return fmt.Errorf("failed to resolve signature configuration for the bundle identifier \"%v\"", bundleID)
	}

	cfg = append(cfg, api.TargetSignatureConfig{
		TargetName: t,
		Config:     sc,
	})

	return nil
}

func (s signatureService) configureDependencies(nt pbx.NativeTarget, f func(string) error) error {
	// Do the native target has native depdendencies
	for _, dp := range nt.Dependencies {
		switch dp.ProductType {
		case pbx.Application, pbx.TvExtension:
			if err := f(dp.Name); err != nil {
				return err
			}
		}
	}

	return nil
}
