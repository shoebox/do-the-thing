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

type ListService interface {
	List() ([]*Install, error)
}

// Install xcode installation definition
type Install struct {
	Path          string
	BundleVersion string
	Version       string
}

// ListService Service to retrieve the list of xcode installation on the sytem
type XCodeListService struct {
	exec util.Exec
	file util.FileService
}

// New create a new instance of the service
func NewXCodeListService(exec util.Exec, file util.FileService) ListService {
	return XCodeListService{exec: exec, file: util.IoUtilFileService{}}
}

// List return all system XCode installation
func (s XCodeListService) List() ([]*Install, error) {
	data, err := s.spotlightSearch()
	if err != nil {
		return nil, err
	}

	list, err := s.parseSpotlightSearchResult(bytes.NewReader(data))
	return list, err
}

func (s XCodeListService) spotlightSearch() ([]byte, error) {
	return s.exec.Exec(MDFIND, XCODE_BUNDLE_IDENTIFIER)
}

func (s XCodeListService) parseSpotlightSearchResult(reader io.Reader) ([]*Install, error) {
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

func (s XCodeListService) validate(path string) (bool, error) {
	return s.file.IsDir(path)
}

func (s XCodeListService) resolveXcode(path string) (*Install, error) {
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