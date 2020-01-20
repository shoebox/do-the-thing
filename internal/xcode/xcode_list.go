package xcode

import (
	"bufio"
	"bytes"
	"dothething/internal/util"
	"io"
	"path/filepath"

	logr "github.com/sirupsen/logrus"
)

const (
	MDFIND                  = "mdfind"
	XCODE_BUNDLE_IDENTIFIER = "kMDItemCFBundleIdentifier == 'com.apple.dt.Xcode'"
	PLIST                   = "/Contents/Info.plist"
)

// Install xcode installation definition
type Install struct {
	Path          string
	BundleVersion string
	Version       string
}

// ListService Service to retrieve the list of xcode installation on the sytem
type ListService struct {
	exec util.Exec
	file util.FileService
}

// New create a new instance of the service
func New() ListService {
	return ListService{exec: util.OsExec{}, file: util.IoUtilFileService{}}
}

// List return all system XCode installation
func (s ListService) List() ([]*Install, error) {
	data, err := s.spotlightSearch()
	if err != nil {
		return nil, err
	}

	list, err := s.parseSpotlightSearchResult(bytes.NewReader(data))
	logr.Println(list, err)
	return list, err
}

func (s ListService) spotlightSearch() ([]byte, error) {
	return s.exec.Exec(MDFIND, XCODE_BUNDLE_IDENTIFIER)
}

func (s ListService) parseSpotlightSearchResult(reader io.Reader) ([]*Install, error) {
	result := []*Install{}

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		path := scanner.Text()
		valid, err := s.validate(path)
		if err != nil || !valid {
			logr.Error(err)
			continue
		}

		x, err := s.resolveXcode(path)
		if err != nil {
			logr.Error(err)
			continue
		}
		result = append(result, x)
	}

	return result, nil
}

func (s ListService) validate(path string) (bool, error) {
	return s.file.IsDir(path)
}

func (s ListService) resolveXcode(path string) (*Install, error) {
	abs, err := filepath.Abs(path + PLIST)
	if err != nil {
		return nil, err
	}

	info := infoPlist{}
	file, err := s.file.OpenAndReadFileContent(abs)
	if err != nil {
		return nil, err
	}

	err = util.DecodeFile(bytes.NewReader(file), &info)
	if err != nil {
		return nil, err
	}

	return &Install{Path: path, Version: info.Version, BundleVersion: info.BundleVersion}, nil
}

type infoPlist struct {
	BundleVersion string `plist:"CFBundleVersion"`
	Version       string `plist:"CFBundleShortVersionString"`
}
