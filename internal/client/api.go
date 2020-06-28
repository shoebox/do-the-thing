package client

import (
	"dothething/internal/action"
	"dothething/internal/config"
	"dothething/internal/destination"
	"dothething/internal/keychain"
	"dothething/internal/signature"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/project"
	"io/ioutil"
)

// API definition interface. This is the central point to retrieve the instance of the diferent
// services used
type API interface {
	ActionBuild() action.ActionBuild
	ActionArchive() action.ActionArchive
	ActionRunTest() action.ActionRunTest
	CertificateService() signature.CertificateService
	DestinationService() destination.Service
	FileService() util.FileService
	KeyChainService() keychain.KeyChain
	ProvisioningService() signature.ProvisioningService
	Signature() signature.Service
	XCodeProjectService() project.ProjectService
	XCodeListService() xcode.ListService
	XCodeSelectService() xcode.SelectService
}

type api struct {
	actionArchive       action.ActionArchive
	actionBuild         action.ActionBuild
	actionRunTest       action.ActionRunTest
	certificateService  signature.CertificateService
	config              config.Config
	destinationService  destination.Service
	executor            util.Executor
	fileUtil            util.FileService
	keychainService     keychain.KeyChain
	listService         xcode.ListService
	projectService      project.ProjectService
	provisioningService signature.ProvisioningService
	runTest             action.ActionRunTest
	selectService       xcode.SelectService
	signatureResolver   signature.Resolver
	signatureService    signature.Service
	xcodeService        xcode.BuildService
}

func (a api) ActionBuild() action.ActionBuild                    { return a.actionBuild }
func (a api) ActionArchive() action.ActionArchive                { return a.actionArchive }
func (a api) ActionRunTest() action.ActionRunTest                { return a.actionRunTest }
func (a api) CertificateService() signature.CertificateService   { return a.certificateService }
func (a api) DestinationService() destination.Service            { return a.destinationService }
func (a api) KeyChainService() keychain.KeyChain                 { return a.keychainService }
func (a api) FileService() util.FileService                      { return a.fileUtil }
func (a api) ProvisioningService() signature.ProvisioningService { return a.provisioningService }
func (a api) Signature() signature.Service                       { return a.signatureService }
func (a api) XCodeListService() xcode.ListService                { return a.listService }
func (a api) XCodeProjectService() project.ProjectService        { return a.projectService }
func (a api) XCodeSelectService() xcode.SelectService            { return a.selectService }

// NewAPIClient create a new instance of the client
func NewAPIClient(config config.Config) (API, error) {
	var err error
	executor := util.NewExecutor()
	fileService := util.NewFileService()
	res := api{executor: executor, fileUtil: fileService}

	//
	res.xcodeService = xcode.NewService(executor, config.Path)
	res.projectService = project.NewProjectService(ioutil.ReadFile, res.xcodeService, executor)

	//
	if res.keychainService, err = keychain.NewKeyChain(executor); err != nil {
		return res, err
	}

	res.actionBuild = action.NewBuild(res.xcodeService, executor)
	res.actionArchive = action.NewArchive(res.xcodeService, executor)
	res.actionRunTest = action.NewActionRun(res.xcodeService, executor)
	res.certificateService = signature.NewCertificateService(config, fileService)
	res.destinationService = destination.NewDestinationService(res.xcodeService, executor)
	res.listService = xcode.NewXCodeListService(executor, fileService)
	res.provisioningService = signature.NewProvisioningService(executor, fileService)
	res.selectService = xcode.NewSelectService(res.listService, executor)
	res.signatureResolver = signature.NewResolver(res.certificateService, res.provisioningService)

	res.signatureService = signature.New(res.projectService, res.signatureResolver)

	return res, err
}
