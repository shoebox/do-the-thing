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
	"log"
	"sync"
)

type apiDef struct {
	config      *api.Config
	executor    api.Executor
	pathService api.PathService
}

var (
	once        sync.Once
	apiKeyChain api.KeyChain
)
var API apiDef = apiDef{executor: util.NewExecutor()}

func (a apiDef) ActionArchive() api.Action { return action.NewArchive(a, a.config) }
func (a apiDef) ActionBuild() api.Action   { return action.NewBuild(a, a.config) }
func (a apiDef) ActionRunTest() api.Action { return action.NewActionRun(a, a.config) }
func (a apiDef) CertificateService() api.CertificateService {
	return signature.NewCertificateService(a, a.config)
}
func (a apiDef) Config() *api.Config { return a.config }
func (a apiDef) DestinationService() api.DestinationService {
	return destination.NewDestinationService(a)
}
func (a apiDef) Exec() api.Executor           { return a.executor }
func (a apiDef) FileService() api.FileService { return util.NewFileService() }
func (a apiDef) KeyChainService() api.KeyChain {
	if apiKeyChain == nil {
		var err error
		apiKeyChain, err = keychain.NewKeyChain(a)
		if err != nil {
			log.Panic(err)
		}
	}

	return apiKeyChain
}

func (a apiDef) ProvisioningService() api.ProvisioningService {
	return signature.NewProvisioningService(a)
}

func (a apiDef) PListBuddyService() api.PListBuddyService { return util.NewPListBuddy(a, a.config) }
func (a apiDef) PathService() api.PathService             { return a.pathService }
func (a apiDef) SignatureResolver() api.SignatureResolver { return signature.NewResolver(a, a.config) }
func (a apiDef) SignatureService() api.SignatureService   { return signature.New(a, a.config) }
func (a apiDef) XCodeBuildService() api.BuildService      { return xcode.NewService(a, a.config) }
func (a apiDef) XCodeListService() api.ListService        { return xcode.NewXCodeListService(a) }
func (a apiDef) XCodeProjectService() api.ProjectService  { return project.NewProjectService(a) }
func (a apiDef) XCodeSelectService() api.SelectService    { return xcode.NewSelectService(a) }

// NewAPIClient create a new instance of the client
func NewAPIClient(config *api.Config) (api.API, error) {
	res := apiDef{config: config, executor: util.NewExecutor()}
	res.pathService = path.NewPathService(res)
	return res, nil
}
