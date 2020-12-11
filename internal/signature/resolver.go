package signature

import (
	"bytes"
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode/pbx"
	"errors"
	"fmt"
	"strings"
)

var ErrNoMatchFound = errors.New("No match found")

// NewResolver creates a new instance of the signature resolver to be use to find the
// right signature configuration for the provided configuration (aka pair of certificate and
// provisioning)
func NewResolver(api api.API, cfg *api.Config) api.SignatureResolver {
	return signatureResolver{api, cfg}
}

// signatureResolver is the implementation of the SignatureResolver interface
type signatureResolver struct {
	api.API
	*api.Config
}

// Resolve will to try to resolve and match of provisioning profile and certficiate aginst the
// provided project configuration
func (r signatureResolver) Resolve(
	ctx context.Context,
	bundleIdentifier string,
	platform pbx.PBXProductType,
) (*api.SignatureConfiguration, error) {
	var err error
	var res api.SignatureConfiguration
	fmt.Println("Resolve")

	// Matching the right provisioning file for the project bundle identifier configuration
	res.ProvisioningProfile, err = r.resolveProvisioningFileFor(ctx,
		r.Config.CodeSignOption.Path,
		bundleIdentifier,
		platform)
	if err != nil {
		return &res, err
	}

	// The provisioning public key to match on
	provisioningPublicKey := res.ProvisioningProfile.Certificates[0].Raw

	// We iterate on all certificates found in the path
	if res.Cert, err = r.findMatchingCert(
		r.API.CertificateService().ResolveInFolder(ctx, r.Config.CodeSignOption.Path),
		provisioningPublicKey,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r signatureResolver) findMatchingCert(certs []*api.P12Certificate, pc []byte) (*api.P12Certificate, error) {
	for _, c := range certs {
		if bytes.Compare(c.Raw, pc) == 0 {
			return c, nil
		}
	}
	return nil, errors.New("Could not find a matching certificate")
}

// resolveProvisioningFileFor will try to resolve a provisioning for the provided configuration
func (r signatureResolver) resolveProvisioningFileFor(
	ctx context.Context,
	path string,
	bundleIdentifier string,
	platform pbx.PBXProductType,
) (*api.ProvisioningProfile, error) {
	found, pp := r.findFor(r.API.ProvisioningService().ResolveProvisioningFilesInFolder(ctx, path),
		bundleIdentifier,
		platform)

	// We have not found a match, raising an error
	if !found {
		fmt.Println("match not found for bundleIdentifier", bundleIdentifier)
		return nil, errors.New("ProvisioningProfile not found")
	}

	return pp, nil
}

func (r signatureResolver) findFor(
	pps []*api.ProvisioningProfile,
	bundleIdentifier string,
	platform pbx.PBXProductType,
) (bool, *api.ProvisioningProfile) {
	// We then iterate on the list to find a match against the project bundle identifier
	for _, pp := range pps {
		if !contains(pp.Platform, platform) {
			continue
		}

		// Do we have a bundle identifier match
		if pp.BundleIdentifier == "*" || pp.BundleIdentifier == bundleIdentifier {
			return true, pp
		}
	}

	return false, nil
}

func contains(a []string, v pbx.PBXProductType) bool {
	var res bool
	for _, s := range a {
		switch v {
		case pbx.TvExtension:
			res = strings.EqualFold(s, "tvOS")

		case pbx.Application:
			res = strings.EqualFold(s, "iOS") || strings.EqualFold(s, "tvOS")
		}

		if res {
			break
		}
	}

	return res
}
