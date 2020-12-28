// +build ignore

package api

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/stretchr/testify/mock"
)

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) OpenAndReadFileContent(abs string) ([]byte, error) {
	c := m.Called(abs)
	return []byte(c.String(0)), c.Error(1)
}

func (m *MockFileService) Open(path string) (io.ReadCloser, error) {
	c := m.Called(path)
	return c.Get(0).(io.ReadCloser), c.Error(1)
}

func (m *MockFileService) IsDir(path string) (bool, error) {
	c := m.Called(path)
	return c.Get(0).(bool), c.Error(1)
}

func (m *MockFileService) Walk(
	ctx context.Context,
	root string,
	isValid func(info os.FileInfo) bool,
	file chan string, wg *sync.WaitGroup,
) error {
	c := m.Called(ctx, root, isValid, file)
	return c.Error(0)
}
