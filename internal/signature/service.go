package signature

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode/pbx"
)

type service struct {
	api.API
	*api.Config
}

func New(api api.API, cfg *api.Config) api.SignatureService {
	return service{api, cfg}
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
	// fmt.Printf(" >>> %#v %v\n", sc, err)

	*res = append(*res, api.TargetSignatureConfig{
		TargetName: t,
		Config:     sc,
	})

	/*
		// Importing the certificate
		if err := s.API.KeyChainService().ImportCertificate(
			ctx,
			sc.Cert.FilePath,
			s.Config.CodeSignOption.CertificatePassword,
		); err != nil {
			return res, err
		}
		fmt.Println("Importing")

		// Installing certificate
		if err := s.installProvisioning(sc.ProvisioningProfile); err != nil {
			return res, err
		}

		res = []string{
			fmt.Sprintf("CODE_SIGN_IDENTITY=%v", sc.Cert.Subject.CommonName),
			"CODE_SIGN_STYLE=Manual",
			fmt.Sprintf("DEVELOPMENT_TEAM=%v", sc.ProvisioningProfile.Entitlements.TeamID),
			fmt.Sprintf("PROVISIONING_PROFILE_SPECIFIER=%v", sc.ProvisioningProfile.UUID),
		}
	*/

	return nil
}

func (s service) Run(ctx context.Context, project api.Project) ([]api.TargetSignatureConfig, error) {
	var res []api.TargetSignatureConfig

	// Found configuration, installing it into a temporary keychain
	if err := s.API.KeyChainService().Create(ctx, "dothething"); err != nil {
		return res, err
	}

	err := s.forTarget(ctx, s.Config.Target, project, &res)
	return res, err
}

/*
func (s api.TargetSignatureConfig) installProvisioning() error {
	input, err := ioutil.ReadFile(s.Config.ProvisioningProfile.FilePath)
	if err != nil {
		return err
	}

		// Retrieving the user home directory
		dir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		// Formatting the provisioning path
		fn := fmt.Sprintf("%v/Library/MobileDevice/Provisioning Profiles/%v.mobileprovision",
			dir,
			p.UUID)

		// Writing the file
		return ioutil.WriteFile(fn, input, os.ModePerm)
}
*/
