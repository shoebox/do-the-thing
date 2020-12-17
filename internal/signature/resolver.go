package signature

import (
	"bytes"
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode/pbx"
	"errors"
	"sort"
	"strings"
)

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

	// resolving the candidates to match against
	candidates := r.API.
		ProvisioningService().
		ResolveProvisioningFilesInFolder(ctx, r.Config.CodeSignOption.Path)

	// Matching the right provisioning file for the project bundle identifier configuration
	if res.ProvisioningProfile, err = r.resolveProvisioningFileFor(
		ctx,
		candidates,
		bundleIdentifier,
		platform,
	); err != nil {
		return nil, err
	}

	// The provisioning public key to match on
	provisioningPublicKey := res.ProvisioningProfile.Certificates[0].Raw

	// We iterate on all certificates found in the path
	certs := r.API.CertificateService().ResolveInFolder(ctx, r.Config.CodeSignOption.Path)

	// And we try to find a matching certificate to the provisioning profile
	if res.Cert, err = r.findMatchingCert(certs, provisioningPublicKey); err != nil {
		return nil, NewSignatureError(err, ErrorCertificateResolution)
	}

	return &res, nil
}

// findMatchingCert will check if a matching certificate can be found into the list
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
	candidates []*api.ProvisioningProfile,
	bundleIdentifier string,
	platform pbx.PBXProductType,
) (*api.ProvisioningProfile, error) {

	// resolving candidate
	found, pp := r.findFor(candidates, bundleIdentifier, platform)

	// we have not found a match, raising an error
	if !found {
		return nil, NewSignatureError(nil, ErrorProvisioningProfileResolution)
	}

	return pp, nil
}

// sortBundleIdentifiers will sort the bundle identifier by expiration date, and wildcards last
func sortBundleIdentifiers(pps []*api.ProvisioningProfile) {
	sort.SliceStable(pps, func(i, j int) bool {
		ppi := pps[i]
		ppj := pps[j]
		// If both bundle identifier are a wildcard, we sort by expiration date
		if ppi.BundleIdentifier == "*" && ppj.BundleIdentifier == "*" {
			return ppj.ExpirationDate.After(ppi.ExpirationDate)
		}

		if ppi.BundleIdentifier == "*" {
			return false
		}

		if ppj.BundleIdentifier == "*" {
			return true
		}

		return ppj.ExpirationDate.After(ppi.ExpirationDate)
	})
}

// findFor will try to resolve a matching provisioning profile into the provided list for the
// required bundle identifier and platform
func (r signatureResolver) findFor(
	pps []*api.ProvisioningProfile,
	bundleIdentifier string,
	platform pbx.PBXProductType,
) (bool, *api.ProvisioningProfile) {
	sortBundleIdentifiers(pps)

	// We then iterate on the list to find a match against the project bundle identifier
	for _, pp := range pps {
		if !contains(pp.Platform, platform) {
			continue
		}

		// Wildcard
		if pp.BundleIdentifier == "*" {
			return true, pp
		}

		// Do we have a bundle identifier match
		if pp.BundleIdentifier == bundleIdentifier {
			return true, pp
		}

		// Wildcard domains
		if strings.HasSuffix(pp.BundleIdentifier, "*") {
			if strings.HasPrefix(bundleIdentifier, strings.TrimSuffix(pp.BundleIdentifier, ".*")) {
				return true, pp
			}
		}
	}

	return false, nil
}

// contains will check that the PBX product type is contained into the provisioning profile capabilities
func contains(a []string, v pbx.PBXProductType) bool {
	// TODO: Adding support of watches...
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
