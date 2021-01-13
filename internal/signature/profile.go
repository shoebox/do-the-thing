package signature

import (
	"bytes"
	"context"
	"crypto/x509"
	"dothething/internal/api"
	"dothething/internal/util"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"go.mozilla.org/pkcs7"
)

var (
	// ErrorFailedToDecode the decode of the provisioning profile failed
	ErrorFailedToDecode = errors.New("Failed to decode the provisioning file")

	// ErrorParsingPublicKey the parsing of the public key contained in the provisioning pofile failed
	ErrorParsingPublicKey = errors.New("Failed to parse the provisioning file certificate")
)

// provisioningService implement the ProvisioningService interface
type provisioningService struct {
	*api.API
}

// NewProvisioningService create a new instance of the provisioning service
func NewProvisioningService(api *api.API) api.ProvisioningService {
	return provisioningService{api}
}

// Decode will decode the provisioning at the designated filepath
func (p provisioningService) Decode(ctx context.Context, r io.Reader) (api.ProvisioningProfile, error) {
	var pp api.ProvisioningProfile

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
	pp.BundleIdentifier = strings.TrimSpace(strings.TrimPrefix(pp.Entitlements.AppID,
		fmt.Sprintf("%s.", pp.Entitlements.TeamID)))

	return pp, err
}

func (p provisioningService) Install(pp *api.ProvisioningProfile) error {
	log.Info().
		Str("ID", pp.UUID).
		Str("BundleIdentifier", pp.BundleIdentifier).
		Msg("Installing provisioning")

	input, err := ioutil.ReadFile(pp.FilePath)
	if err != nil {
		return err
	}

	// Retrieving the user home directory
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Formatting the provisioning path
	fn := fmt.Sprintf("%v/Library/MobileDevice/Provisioning Profiles/%v.mobileprovision",
		dir,
		pp.UUID)

	// Writing the file
	return ioutil.WriteFile(fn, input, os.ModePerm)
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

func (p provisioningService) decodeRawProvisioning(
	ctx context.Context,
	filepath string,
	reader io.ReadCloser,
) (*api.ProvisioningProfile, error) {
	defer reader.Close()

	// Which try to decode the candidate provisioning file
	dpp, err := p.Decode(ctx, reader)

	if err != nil {
		// Here we do not return the error, like we do not want the parent errgroup to fail
		return nil, err
	}

	dpp.FilePath = filepath

	return &dpp, nil
}

// readProvisioningFile Read the provided file at path and try to decode it as provisioning
func (p provisioningService) readProvisioningFile(
	ctx context.Context,
	path string,
) (*api.ProvisioningProfile, error) {
	// Open the file to a reader
	f, err := p.API.FileService.Open(path)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	// Read the content of the file
	return p.decodeRawProvisioning(ctx, path, f)
}

func (p provisioningService) decodingWorker(
	id int,
	wg *sync.WaitGroup,
	ctx context.Context,
	paths <-chan string,
	res *[]*api.ProvisioningProfile,
) {
	wg.Add(1)
	for value := range paths {
		dpp, err := p.readProvisioningFile(ctx, value)
		if err != nil {
			fmt.Println("error while decoding", err)
		} else {
			*res = append(*res, dpp)
		}
	}

	wg.Done()
}

// ResolveProvisioningFilesInFolder walk the provided root path and resolve all provisioning
// profiles contained into it
func (p provisioningService) ResolveProvisioningFilesInFolder(
	ctx context.Context,
	root string,
) []*api.ProvisioningProfile {

	paths := make(chan string)
	var res []*api.ProvisioningProfile

	var wg sync.WaitGroup
	// Increment waitgroup counter and create go routines
	for i := 0; i < 8; i++ {
		go p.decodingWorker(i, &wg, ctx, paths, &res)
	}

	err := p.API.FileService.Walk(ctx, root, isProvisioningFile, paths, &wg)

	if err != nil {
		fmt.Println("err ", err)
	}
	wg.Wait()

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
