package signature

import (
	"bytes"
	"context"
	"dothething/internal/api"
	"errors"
	"fmt"
)

var ErrNoMatchFound = errors.New("No match found")

// NewResolver creates a new instance of the signature resolver to be use to find the
// right signature configuration for the provided configuration (aka pair of certificate and
// provisioning)
func NewResolver(api api.API) api.SignatureResolver {
	return signatureResolver{api}
}

// signatureResolver is the implementation of the SignatureResolver interface
type signatureResolver struct {
	api.API
}

// Resolve will to try to resolve and match of provisioning profile and certficiate aginst the
// provided project configuration
func (r signatureResolver) Resolve(ctx context.Context,
	// config config.Config,
	path string,
	bundleIdentifier string) (api.SignatureConfiguration, error) {

	var err error
	var res api.SignatureConfiguration

	// Matching the right provisioning file for the project bundle identifier configuration
	res.ProvisioningProfile, err = r.resolveProvisioningFileFor(ctx,
		path,
		bundleIdentifier)
	if err != nil {
		return res, err
	}

	// And trying find a matching certificate to pair with the bundle identifier
	certs := r.API.CertificateService().ResolveInFolder(ctx, path)
	var found bool

	// The provisioning public key to match on
	provisioningPublicKey := res.ProvisioningProfile.Certificates[0].Raw

	// We iterate on all certificates found in the path
	for _, c := range certs {
		fmt.Println(string(c.Raw))
		// We check if the certificate public key is matching the provisioning's
		if bytes.Compare(c.Raw, provisioningPublicKey) == 0 {
			// If yes we created the new pair object with those.
			res.Cert = c
			found = true
		}
	}

	if !found {
		return res, errors.New("Could not find a matching certificate")
	}

	return res, nil
}

// resolveProvisioningFileFor will try to resolve a provisioning for the provided configuration
func (r signatureResolver) resolveProvisioningFileFor(ctx context.Context,
	path string,
	bundleIdentifier string) (api.ProvisioningProfile, error) {

	var res api.ProvisioningProfile
	var err error
	var found bool

	// Resolving all provisining in the folder
	pps := r.API.ProvisioningService().ResolveProvisioningFilesInFolder(ctx, path)

	// We then iterate on the list to find a match against the project bundle identifier
	for _, pp := range pps {
		// Do we have a bundle identifier match
		if pp.BundleIdentifier == "*" || pp.BundleIdentifier == bundleIdentifier {
			found = true
			res = pp
			break
		}
	}

	// We have not found a match, raising an error
	if !found {
		return res, errors.New("ProvisioningProfile not found")
	}

	return res, err
}
