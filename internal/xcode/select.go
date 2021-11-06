package xcode

import (
	"context"
	"dothething/internal/api"
	"fmt"
	"sort"

	"github.com/blang/semver"
	"github.com/rs/zerolog/log"
)

var (
	// ErrMatchNotFound No match found on the system
	ErrMatchNotFound = "XCode match not found"

	// ErrInvalidVersion the required vesion format is invalid
	ErrInvalidVersion = "Invalid version"

	ErrParsing = "Failed to parse required version"
)

// selectService The XCode version selection service struct
type selectService struct{ *api.API }

// NewSelectService create a new instance of the xcode selector service
func NewSelectService(api *api.API) api.SelectService {
	return selectService{api}
}

// Find allow to resolve a XCode install by required vesion
func (s selectService) Find(ctx context.Context) (*api.Install, error) {
	log.Info().
		Str("Requirement", s.API.Config.XCodeVersion).
		Msg("Finding XCode installation")

	r, err := semver.ParseRange(s.API.Config.XCodeVersion)
	if err != nil {
		fmt.Println("err :::", err)
		return nil, fmt.Errorf("%v (%v)", ErrParsing, err)
	}


	// Find a equal match
	target, err := s.findMatch(ctx, r, s.isMatchingRequirement)

	// In case of no match found
	if target == nil || err != nil {
		return nil, fmt.Errorf("%v (%v)", ErrMatchNotFound, err)
	}

	return target, nil
}

func (s *selectService) isMatchingRequirement(i *api.Install, r semver.Range) (bool, error) {
	v, err := semver.ParseTolerant(i.Version)
	if err != nil {
		return false, err
	}

	return r(v), nil
}

func (s *selectService) findMatch(
	ctx context.Context,
	r semver.Range,
	valid func(install *api.Install, r semver.Range) (bool, error),
) (*api.Install, error) {
	// Resolve the list of candidates
	list, err := s.API.XcodeListService.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to list system XCode installations (%v)", err)
	}

	// Iterate on installs
	var installs []*api.Install
	for _, install := range list {
		res, err := valid(install, r)
		fmt.Println("res,", res, err)
		if err != nil {
			log.Err(err)
			continue
		}

		if res {
			installs = append(installs, install)
		}
	}

	// Sort the installs
	sortInstalls(installs)

	// In case of no candidates
	if len(installs) == 0 {
		return nil, fmt.Errorf("%v", ErrMatchNotFound)
	}

	return installs[0], nil
}

func sortInstalls(installs []*api.Install) {
	sort.Slice(installs, func(i, j int) bool {
		return compareInstall(installs[i], installs[j])
	})
}

func compareInstall(i1 *api.Install, i2 *api.Install) bool {
	v1, err := semver.Parse(i1.Version)
	if err != nil {
		return false
	}

	v2, err := semver.Parse(i2.Version)
	if err != nil {
		return false
	}

	return v1.GT(v2)
}
