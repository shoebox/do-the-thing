package signature

import (
	"bytes"
	"context"
	"crypto/x509"
	"dothething/internal/util"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"go.mozilla.org/pkcs7"
)

// ProvisioningProfile type definition
type ProvisioningProfile struct {
	BundleIdentifier string
	Certificates     []*x509.Certificate
	Entitlements     Entitlements `plist:"Entitlements"`
	ExpirationDate   time.Time    `plist:"ExpirationDate"`
	Name             string       `plist:"Name"`
	Platform         []string     `plist:"Platform"`
	RawCertificates  [][]byte     `plist:"DeveloperCertificates"`
	TeamName         string       `plist:"TeamName"`
	UUID             string       `plist:"UUID"`
}

// Entitlements provisioning entitlements definition
type Entitlements struct {
	AccessGroup string `json:"keychain-access-groups"`
	Aps         string `json:"aps-environment"`
	AppID       string `plist:"application-identifier"`
	TeamID      string `plist:"com.apple.developer.team-identifier"`
}

var (
	// ErrorFailedToDecode the decode of the provisioning profile failed
	ErrorFailedToDecode = errors.New("Failed to decode the provisioning file")

	// ErrorParsingPublicKey the parsing of the public key contained in the provisioning pofile failed
	ErrorParsingPublicKey = errors.New("Failed to parse the provisioning file certificate")
)

// ProvisioningService interface to describe the provisioning service method
type ProvisioningService interface {
	Decode(ctx context.Context, r io.Reader) (ProvisioningProfile, error)
	ResolveProvisioningFilesInFolder(ctx context.Context, root string) []ProvisioningProfile
}

// provisioningService implement the ProvisioningService interface
type provisioningService struct {
	util.FileService
	util.Executor
}

// NewProvisioningService create a new instance of the provisioning service
func NewProvisioningService(e util.Executor, f util.FileService) ProvisioningService {
	return provisioningService{Executor: e, FileService: f}
}

// Decode will decode the provisioning at the designated filepath
func (p provisioningService) Decode(ctx context.Context, r io.Reader) (ProvisioningProfile, error) {
	var pp ProvisioningProfile

	// First we decode the provisioning at path
	data, err := p.decodeProvisioning(ctx, r)
	if err != nil {
		return pp, ErrorParsingPublicKey
	}

	// We parse the provisioning plist file content, and unmarshal it
	if err := util.DecodeFile(bytes.NewReader(data), &pp); err != nil {
		return pp, ErrorFailedToDecode
	}

	// Parse raw x509 Certificates
	pp.Certificates, err = parseRawX509Certificates(pp.RawCertificates)
	if err != nil {
		return pp, err
	}

	// For more convenience compute the bundle identifier without the teamID prefix.
	pp.BundleIdentifier = strings.TrimPrefix(pp.Entitlements.AppID,
		fmt.Sprintf("%s.", pp.Entitlements.TeamID))

	return pp, err
}

// isProvisioningFile check if the candidate is a valid provisioning file by testing it's
// type and file extension
func isProvisioningFile(info os.FileInfo) bool {
	return info.Mode().IsRegular() &&
		!info.IsDir() &&
		strings.HasSuffix(info.Name(), ".mobileprovision")
}

// walkOnPath validate if the provided path is provisioning file and report the result
// to be paired with the filepath.Walk call in ResolveProvisioningFilesInFolder
func (p provisioningService) walkOnPath(
	ctx context.Context,
	path string,
	cpath chan string,
	info os.FileInfo,
) error {

	// Check if the file is a provisioning file
	if !isProvisioningFile(info) {
		return nil
	}

	// Select result
	select {
	case cpath <- path: // Add the result path to the channel
	case <-ctx.Done(): // Handle context cancelation
		return ctx.Err()
	}
	return nil
}

func (p provisioningService) decodeRawProvisioning(ctx context.Context,
	reader io.Reader,
	cprovisioning chan ProvisioningProfile) error {

	// Which try to decode the candidate provisioning file
	dpp, err := p.Decode(ctx, reader)
	if err != nil {
		return err
	}

	//Use a select statement to exit out if context expires
	select {
	case <-ctx.Done(): // Handle context cancelation
		return nil
	case cprovisioning <- dpp: // Add the decoded provisioning to the channel
		return nil
	}
}

// readProvisioningFile Read the provided file at path and try to decode it as provisioning
func (p provisioningService) readProvisioningFile(
	ctx context.Context,
	path string,
	cprovisioning chan ProvisioningProfile,
) error {
	// Open the file to a reader
	f, err := p.FileService.Open(path)
	if err != nil {
		return err
	}

	// Close the file
	defer f.Close()

	// Read the content of the file
	return p.decodeRawProvisioning(ctx, f, cprovisioning)
}

// ResolveProvisioningFilesInFolder walk the provided root path and resolve all provisioning
// profiles contained into it
func (p provisioningService) ResolveProvisioningFilesInFolder(
	ctx context.Context,
	root string,
) []ProvisioningProfile {
	res := []ProvisioningProfile{}

	// Channel of provisionings
	provisionings := make(chan ProvisioningProfile)
	errgroup := p.FileService.Walk(
		ctx,
		root,
		isProvisioningFile,
		// For each candidate certificate file
		func(ctx context.Context, path string) error {
			return p.readProvisioningFile(ctx, path, provisionings)
		})

	go func() {
		// Waiting for the goroutines to complete
		errgroup.Wait()

		// Closing the channel
		close(provisionings)
	}()

	// Convert the channel result to a slice
	for pp := range provisionings {
		res = append(res, pp)
	}

	// Check whether any of the goroutines failed. Since the error group is accumulating the
	// errors, we don't need to send them (or check for them) in the individual
	// results sent on the channel.
	if err := errgroup.Wait(); err != nil {
		return res
	}

	return res
}

// decodeProvisioning is using the security API to decode the provisioning file
func (p provisioningService) decodeProvisioning(ctx context.Context, r io.Reader) ([]byte, error) {
	var res []byte

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return res, err
	}

	// Decrypt the DMS message encrypted (DMS is based on PKCS#7)
	// which is the equivalent of using the security cms toolkit
	p7, err := pkcs7.Parse(b)
	if err != nil {
		return res, err
	}

	// Return message content
	return p7.Content, nil
}

// parseRawX509Certificates will parse the raw certificate slice
func parseRawX509Certificates(raw [][]byte) ([]*x509.Certificate, error) {
	res := []*x509.Certificate{}

	// Iterate on all raw certificates
	for _, data := range raw {
		// To parse them
		k, err := x509.ParseCertificate(data)
		if err != nil {
			return nil, ErrorParsingPublicKey
		}

		// And append the parse certificate to the result array
		res = append(res, k)
	}

	return res, nil
}
