package signature

import (
	"bytes"
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

func DecodeProvisioningProfile(filePath string, exec util.Exec) (*ProvisioningProfile, error) {
	var pp ProvisioningProfile
	data, err := decodeProvisioning(filePath, exec)
	if err != nil {
		return nil, err
	}

	// Parse plist file
	err = util.DecodeFile(bytes.NewReader(data), &pp)
	if err != nil {
		return nil, ErrorFailedToDecode
	}

	// Parse raw x509 Certificates
	pp.Certificates, err = parseRawX509Certificates(pp.RawCertificates)
	if err != nil {
		return nil, err
	}

	return &pp, err
}

func decodeProvisioning(filePath string, exec util.Exec) ([]byte, error) {
	// Invoke security tool to decode the file
	data, err := exec.Exec(nil, Security, Cms, ArgDecodeCMS, ArgInlineFile, filePath)
	if err != nil {
		return nil, ErrorFailedToDecode
	}

	return data, err
}

func parseRawX509Certificates(raw [][]byte) ([]*x509.Certificate, error) {
	certs := []*x509.Certificate{}
	for _, data := range raw {
		k, err := x509.ParseCertificate(data)
		if err != nil {
			return nil, ErrorParsingPublicKey
		}

		certs = append(certs, k)
	}

	return certs, nil
}
