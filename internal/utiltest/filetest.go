package utiltest

import (
	"github.com/stretchr/testify/mock"
)

type MockFileService struct {
	mock.Mock
}

func (f *MockFileService) OpenAndReadFileContent(path string) ([]byte, error) {
	args := f.Called(path)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return []byte(args.String(0)), nil
}

func (f *MockFileService) IsDir(path string) (bool, error) {
	args := f.Called(path)

	if args.Error(1) != nil {
		return false, args.Error(1)
	}

	return args.Bool(0), nil
}
