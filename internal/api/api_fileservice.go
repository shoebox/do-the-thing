package api

import (
	"context"
	"io"
	"os"
	"sync"
)

type FileService interface {
	OpenAndReadFileContent(abs string) ([]byte, error)
	Open(path string) (io.ReadCloser, error)
	IsDir(path string) (bool, error)
	Walk(ctx context.Context,
		root string,
		isValid func(info os.FileInfo) bool,
		file chan string,
		wg *sync.WaitGroup,
	) error
	// handlePath func(ctx context.Context, path string) error) *errgroup.Group
}
