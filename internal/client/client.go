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
func NewAPIClient() (*api.API, error) {
	a := api.API{Config: &api.Config{}}
	a.Exec = util.NewExecutor(&a)

	// keychain
	k, err := keychain.NewKeyChain(&a)
	if err != nil {
		return nil, err
	}
	a.KeyChain = k

	a.ActionArchive = action.NewArchive(&a)
	a.ActionBuild = action.NewBuild(&a)
	a.ActionPack = action.NewActionPackage(&a)
	a.ActionRunTest = action.NewActionRunTest(&a)

	a.BuildService = xcode.NewService(&a)
	a.CertificateService = signature.NewCertificateService(&a)
	a.DestinationService = destination.NewDestinationService(&a)
	a.ExportOptionService = signature.NewExportOptionsService(&a)
	a.FileService = util.NewFileService()
	a.PathService = path.NewPathService(&a)
	a.PlistBuddyService = util.NewPListBuddy(&a)
	a.ProvisioningService = signature.NewProvisioningService(&a)
	a.SignatureResolver = signature.NewResolver(&a)
	a.SignatureService = signature.NewSignatureService(&a)
	a.XCodeProjectService = project.NewProjectService(&a)
	a.XcodeListService = xcode.NewXCodeListService(&a)
	a.XcodeSelectService = xcode.NewSelectService(&a)
	return &a, nil
}
