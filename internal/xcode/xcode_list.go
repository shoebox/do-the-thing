package xcode

import (
	"bufio"
	"bytes"
	"context"
	"dothething/internal/util"
	"io"
	"path/filepath"

	logr "github.com/sirupsen/logrus"
)

const (
	// MdFind spotlight search executable
	MdFind = "mdfind"

	// XCodeBundleIdentifier What spotlight should look for to identify Xcode installs
	XCodeBundleIdentifier = "kMDItemCFBundleIdentifier == 'com.apple.dt.Xcode'"

	// ContentPListFile path to the Info plist file in to the Xcode app bundle
	ContentPListFile = "/Contents/Info.plist"
)

// ListService basic interface
type ListService interface {
	List(ctx context.Context) ([]*Install, error)
}

// Install xcode installation definition
type Install struct {
	Path          string
	BundleVersion string
	Version       string
}

// listService Service to retrieve the list of xcode installation on the system
type listService struct {
	exec util.Executor
	file util.FileService
}

// New create a new instance of the service
func NewXCodeListService(exec util.Executor, file util.FileService) ListService {
	return listService{exec: exec, file: util.IoUtilFileService{}}
}

// List return all system XCode installation
func (s listService) List(ctx context.Context) ([]*Install, error) {
	data, err := s.spotlightSearch(ctx)
	if err != nil {
		return nil, err
	}

	return s.parseSpotlightSearchResult(bytes.NewReader(data))
}

func (s listService) spotlightSearch(ctx context.Context) ([]byte, error) {
	return s.exec.CommandContext(ctx, MdFind, XCodeBundleIdentifier).Output()
}

func (s listService) parseSpotlightSearchResult(reader io.Reader) ([]*Install, error) {
	result := []*Install{}
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

func (s listService) parseSpotlightEntry(path string) (*Install, error) {
	if valid, err := s.validate(path); err != nil || !valid {
		return nil, err
	}

	return s.resolveXcode(path)
}

func (s listService) validate(path string) (bool, error) {
	return s.file.IsDir(path)
}

func (s listService) resolveXcode(path string) (*Install, error) {
	abs, err := filepath.Abs(path + ContentPListFile)
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
