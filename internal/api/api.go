package api

import (
	"context"
	"crypto/x509"
	"io"
)

type CertificateService interface {
	DecodeCertificate(r io.Reader, password string) (P12Certificate, error)
	ResolveInFolder(ctx context.Context, root string) []P12Certificate
}

// API definition interface. This is the central point to retrieve the instance of the diferent
// services used
type API interface {
	ActionArchive() Action
	ActionBuild() Action
	ActionRunTest() Action
	CertificateService() CertificateService
	DestinationService() DestinationService
	Exec() Executor
	FileService() FileService
	KeyChainService() KeyChain
	ProvisioningService() ProvisioningService
	SignatureResolver() SignatureResolver
	SignatureService() SignatureService
	XCodeBuildService() BuildService
	XCodeListService() ListService
	XCodeProjectService() ProjectService
	XCodeSelectService() SelectService
}

type Action interface {
	Run(ctx context.Context, config Config) error
}

// P12Certificate is a more convenient alias
type P12Certificate struct {
	*x509.Certificate
	FilePath string
}
