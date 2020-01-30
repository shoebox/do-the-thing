package xcode

import (
	"dothething/internal/util"
	"errors"
	"sort"

	"github.com/blang/semver"
	logr "github.com/sirupsen/logrus"
)

var (
	ErrMatchNotFound  = errors.New("XCode match not found")
	ErrInvalidVersion = errors.New("Invalid version")
)

type SelectService interface {
	SelectService(version string) error
}

type XCodeSelectService struct {
	exec util.Exec
	list ListService
}

func NewSelectService(list ListService, exec util.Exec) *XCodeSelectService {
	return &XCodeSelectService{exec: exec, list: list}
}

func (s *XCodeSelectService) SelectVersion(requirement string) (*Install, error) {
	// Should try to parse the required version
	required, err := semver.Parse(requirement)
	if err != nil {
		return nil, ErrInvalidVersion
	}

	// Find a equal match
	target, err := s.findMatch(required, s.isEqualMatch)

	// In case of no match found
	if target == nil || err != nil {
		return nil, ErrMatchNotFound
	}

	return target, nil
}

func (s *XCodeSelectService) isEqualMatch(install *Install, requirement semver.Version) (bool, error) {
	if v, err := semver.Parse(install.Version); err == nil {
		if v.Equals(requirement) {
			return true, nil
		}
	}
	return false, nil
}

func (s *XCodeSelectService) findMatch(requirement semver.Version,
	valid func(install *Install, version semver.Version) (bool, error)) (*Install, error) {

	// Resolve the list of candidates
	list, err := s.list.List()
	if err != nil {
		return nil, err
	}

	// Iterate on installs
	var installs = []*Install{}
	for _, install := range list {
		res, err := valid(install, requirement)
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

func sortInstalls(installs []*Install) {
	sort.Slice(installs, func(i, j int) bool {
		return compareInstall(installs[i], installs[j])
	})
}

func compareInstall(i1 *Install, i2 *Install) bool {
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
