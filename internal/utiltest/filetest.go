package utiltest

import (
	"context"
	"io"
	"os"

	"github.com/stretchr/testify/mock"
	"golang.org/x/sync/errgroup"
)

type MockFileService struct {
	mock.Mock
}

func (f *MockFileService) Walk(ctx context.Context,
	root string,
	isValid func(info os.FileInfo) bool,
	handlePath func(ctx context.Context, path string) error) *errgroup.Group {
	args := f.Called(ctx, root, isValid, handlePath)

	a := args.Get(0)
	g, ok := a.(*errgroup.Group)
	if !ok {
		panic("Wrong type")
	}

	return g
}

func (f *MockFileService) Open(path string) (io.ReadCloser, error) {
	args := f.Called(path)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	a := args.Get(0)
	file, ok := a.(io.ReadCloser)
	if !ok {
		panic("Wrong type")
	}
	return file, nil
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
