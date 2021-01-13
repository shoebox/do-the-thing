package client

import (
	"dothething/internal/action"
	"dothething/internal/api"
	"dothething/internal/destination"
	"dothething/internal/keychain"
	"dothething/internal/path"
	"dothething/internal/signature"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/project"
)

// NewAPIClient create a new instance of the client
func NewAPIClient(config *api.Config) (*api.API, error) {
	res := api.API{Config: config, Exec: util.NewExecutor()}
	res.ActionPack = action.NewActionPackage(&res)
	res.PathService = path.NewPathService(&res)

	// keychain
	k, err := keychain.NewKeyChain(&res)
	if err != nil {
		return nil, err
	}
	res.Config = config
	res.KeyChain = k

	res.ExportOptionService = signature.NewExportOptionsService(&res)
	res.SignatureService = signature.NewSignatureService(&res)
	res.ActionArchive = action.NewArchive(&res)
	res.BuildService = xcode.NewService(&res)
	res.CertificateService = signature.NewCertificateService(&res)
	res.DestinationService = destination.NewDestinationService(&res)
	res.FileService = util.NewFileService()
	res.PlistBuddyService = util.NewPListBuddy(&res)
	res.ProvisioningService = signature.NewProvisioningService(&res)
	res.SignatureResolver = signature.NewResolver(&res)
	res.XcodeListService = xcode.NewXCodeListService(&res)
	res.XCodeProjectService = project.NewProjectService(&res)
	res.XcodeSelectService = xcode.NewSelectService(&res)

	return &res, nil
}
