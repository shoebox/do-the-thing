package xcode

import (
	"context"
	"dothething/internal/api"
	"errors"
	"sort"

	"github.com/blang/semver"
	logr "github.com/sirupsen/logrus"
)

var (
	// ErrMatchNotFound No match found on the system
	ErrMatchNotFound = errors.New("XCode match not found")

	// ErrInvalidVersion the required vesion format is invalid
	ErrInvalidVersion = errors.New("Invalid version")

	ErrParsing = errors.New("Failed to parse required version")
)

// selectService The XCode version selection service struct
type selectService struct{ api.API }

// NewSelectService create a new instance of the xcode selector service
func NewSelectService(api api.API) api.SelectService {
	return selectService{api}
}

// Find allow to resolve a XCode install by required vesion
func (s selectService) Find(ctx context.Context, req string) (*api.Install, error) {
	r, err := semver.ParseRange(req)
	if err != nil {
		return nil, ErrParsing
	}

	// Find a equal match
	target, err := s.findMatch(ctx, r, s.isMatchingRequirement)

	// In case of no match found
	if target == nil || err != nil {
		return nil, ErrMatchNotFound
	}

	return target, nil
}

func (s *selectService) isMatchingRequirement(i *api.Install, r semver.Range) (bool, error) {
	v, err := semver.Parse(i.Version)
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
	list, err := s.API.XCodeListService().List(ctx)
	if err != nil {
		return nil, err
	}

	// Iterate on installs
	var installs []*api.Install
	for _, install := range list {
		res, err := valid(install, r)
		if err != nil {
			logr.Error(err)
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
		return nil, ErrMatchNotFound
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
