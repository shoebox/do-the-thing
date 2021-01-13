package api

import (
	"context"
	"crypto/x509"
	"io"
)

type CertificateService interface {
	DecodeCertificate(r io.Reader, password string) (P12Certificate, error)
	ResolveInFolder(ctx context.Context, root string) []*P12Certificate
}

type API struct {
	ActionArchive       Action
	ActionBuild         Action
	ActionPack          Action
	ActionRun           Action
	ActionRunTest       Action
	BuildService        BuildService
	CertificateService  CertificateService
	Config              *Config
	FileService         FileService
	DestinationService  DestinationService
	Exec                Executor
	ExportOptionService ExportOptionsService
	KeyChain
	PathService         PathService
	PlistBuddyService   PListBuddyService
	ProvisioningService ProvisioningService
	SignatureResolver   SignatureResolver
	SignatureService    SignatureService
	XcodeListService    ListService
	XCodeProjectService ProjectService
	XcodeSelectService  SelectService
}

type Action interface {
	Run(ctx context.Context) error
}

// P12Certificate is a more convenient alias
type P12Certificate struct {
	*x509.Certificate
	FilePath string
}
