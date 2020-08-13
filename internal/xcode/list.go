package xcode

import (
	"bufio"
	"bytes"
	"context"
	"dothething/internal/api"
	"dothething/internal/util"
	"io"
	"path/filepath"

	logr "github.com/sirupsen/logrus"
)

const (
	// MdFind spotlight search executable
	MdFind = "mdfind"

	// XCodeBundleIdentifier What spotlight should look for to identify Xcode installs
	BundleIdentifier = "kMDItemCFBundleIdentifier == 'com.apple.dt.Xcode'"

	// ContentPListFile path to the Info plist file in to the Xcode app bundle
	ContentPListFile = "/Contents/Info.plist"
)

// listService Service to retrieve the list of xcode installation on the system
type listService struct{ api.API }

// NewXCodeListService create a new instance of the service
func NewXCodeListService(api api.API) api.ListService {
	return listService{api}
}

// List return all system XCode installation
func (s listService) List(ctx context.Context) ([]*api.Install, error) {
	data, err := s.spotlightSearch(ctx)
	if err != nil {
		return nil, err
	}

	return s.parseSpotlightSearchResult(bytes.NewReader(data))
}

func (s listService) spotlightSearch(ctx context.Context) ([]byte, error) {
	return s.API.Exec().CommandContext(ctx, MdFind, BundleIdentifier).Output()
}

func (s listService) parseSpotlightSearchResult(reader io.Reader) ([]*api.Install, error) {
	var result []*api.Install
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		path := scanner.Text()
		i, err := s.parseSpotlightEntry(path)
		if err != nil {
			logr.Error(err)
			continue
		}

		result = append(result, i)
	}

	return result, nil
}

func (s listService) parseSpotlightEntry(path string) (*api.Install, error) {
	if valid, err := s.validate(path); err != nil || !valid {
		return nil, err
	}

	return s.resolveXcode(path)
}

func (s listService) validate(path string) (bool, error) {
	return s.API.FileService().IsDir(path)
}

func (s listService) resolveXcode(path string) (*api.Install, error) {
	abs, err := filepath.Abs(path + ContentPListFile)
	if err != nil {
		return nil, err
	}

	info := infoPlist{}

	file, err := s.API.FileService().OpenAndReadFileContent(abs)
	if err != nil {
		return nil, err
	}

	err = util.DecodeFile(bytes.NewReader(file), &info)
	if err != nil {
		return nil, err
	}

	return &api.Install{Path: path, Version: info.Version, BundleVersion: info.BundleVersion}, nil
}

type infoPlist struct {
	BundleVersion string `plist:"CFBundleVersion"`
	Version       string `plist:"CFBundleShortVersionString"`
}
