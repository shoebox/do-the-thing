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
	installs, err := s.list.List()
	if err != nil {
		return nil, err
	}

	required, err := semver.Parse(requirement)
	if err != nil {
		return nil, ErrInvalidVersion
	}

	for _, install := range installs {
		v, err := semver.Make(install.Version)
		if err != nil {
			logr.Error(err)
			continue
		}

		if v.Equals(required) {
			return install, nil
		}
	}

	target, err := s.findMatch(requirement, s.isEqualMatch)

	if target == nil {
		return nil, ErrMatchNotFound
	}

	return target, nil
}

func (s *XCodeSelectService) isEqualMatch(install *Install, version string) (bool, error) {
	r, err := semver.Parse(version)
	if err != nil {
		return false, err
	}

	if v, err := semver.Parse(install.Version); err == nil {
		if v.Equals(r) {
			return true, nil
		}
	}
	return false, nil
}

func (s *XCodeSelectService) findMatch(requirement string,
	valid func(install *Install, version string) (bool, error)) (*Install, error) {
	list, err := s.list.List()
	if err != nil {
		return nil, err
	}

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

	sortInstalls(installs)

	if len(installs) == 0 {
		return nil, ErrMatchNotFound
	}

	return installs[0], nil
}

func sortInstalls(installs []*Install) {
	sort.Slice(installs, func(i, j int) bool {
		version1 := installs[i]
		version2 := installs[j]

		v1, err := semver.Parse(version1.Version)
		if err != nil {
			return false
		}

		v2, err := semver.Parse(version2.Version)
		if err != nil {
			return true
		}

		return v1.GT(v2)
	})
}
