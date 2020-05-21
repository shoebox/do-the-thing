package signature

import (
	"bytes"
	"context"
	"crypto/x509"
	"dothething/internal/util"
	"errors"
)

type ProvisioningProfile struct {
	AppID           string   `plist:"AppIDName"`
	RawCertificates [][]byte `plist:"DeveloperCertificates"`
	Certificates    []*x509.Certificate
	TeamIdentifier  []string `plist:"TeamIdentifier"`
	TeamName        string   `plist:"TeamName"`
	UUID            string   `plist:"UUID"`
}

const (
	Security      string = "security"
	Cms           string = "cms"
	ArgDecodeCMS  string = "-D"
	ArgInlineFile string = "-i"
)

var (
	ErrorFailedToDecode   = errors.New("Failed to decode the provisioning file")
	ErrorParsingPublicKey = errors.New("Failed to parse the provisioning file certificate")
)

type ProvisioningService interface {
	Decode(ctx context.Context, filePath string)
}

type provisioningService struct {
	util.Executor
}

func NewProvisioningService(e util.Executor) provisioningService {
	return provisioningService{Executor: e}
}

func (p provisioningService) Decode(ctx context.Context, filePath string) (ProvisioningProfile, error) {
	var pp ProvisioningProfile
	data, err := p.decodeProvisioning(ctx, filePath)
	if err != nil {
		return pp, err
	}

	// Parse plist file
	if err := util.DecodeFile(bytes.NewReader(data), &pp); err != nil {
		return pp, ErrorFailedToDecode
	}

	// Parse raw x509 Certificates
	pp.Certificates, err = parseRawX509Certificates(pp.RawCertificates)
	if err != nil {
		return pp, err
	}

	return pp, err
}

func (p provisioningService) decodeProvisioning(ctx context.Context, filePath string) ([]byte, error) {
	// Invoke security tool to decode the file
	return p.Executor.CommandContext(ctx,
		Security,
		Cms, ArgDecodeCMS,
		ArgInlineFile, filePath).Output()
}

func parseRawX509Certificates(raw [][]byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for _, data := range raw {
		k, err := x509.ParseCertificate(data)
		if err != nil {
			return nil, ErrorParsingPublicKey
		}

		certs = append(certs, k)
	}

	return certs, nil
}
