package api

import (
	"context"
	"crypto/x509"
	"io"
	"time"
)

type SignatureService interface {
	Run(ctx context.Context) error
}

type TargetSignatureConfig struct {
	TargetName string
	Config     *SignatureConfiguration
}

// ProvisioningService interface to describe the provisioning service method
type ProvisioningService interface {
	Decode(ctx context.Context, r io.Reader) (ProvisioningProfile, error)
	ResolveProvisioningFilesInFolder(ctx context.Context, root string) []*ProvisioningProfile
	Install(p *ProvisioningProfile) error
}

// ProvisioningProfile type definition
type ProvisioningProfile struct {
	BundleIdentifier string
	Certificates     []*x509.Certificate
	Entitlements     Entitlements `plist:"Entitlements"`
	ExpirationDate   time.Time    `plist:"ExpirationDate"`
	Name             string       `plist:"Name"`
	FilePath         string
	Platform         []string `plist:"Platform"`
	RawCertificates  [][]byte `plist:"DeveloperCertificates"`
	TeamName         string   `plist:"TeamName"`
	UUID             string   `plist:"UUID"`
}

// Entitlements provisioning entitlements definition
type Entitlements struct {
	AccessGroup string `json:"keychain-access-groups"`
	Aps         string `json:"aps-environment"`
	AppID       string `plist:"application-identifier"`
	TeamID      string `plist:"com.apple.developer.team-identifier"`
}

// Resolver is the base interface for the signature result
type SignatureResolver interface {
	Resolve(ctx context.Context, bundleIdentifier string, platform string) (*SignatureConfiguration, error)
}

type SignatureConfiguration struct {
	ProvisioningProfile *ProvisioningProfile
	Cert                *P12Certificate
	path                string
}
