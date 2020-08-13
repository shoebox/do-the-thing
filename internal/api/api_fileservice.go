package api

import (
	"context"
	"io"
	"os"

	"golang.org/x/sync/errgroup"
)

type FileService interface {
	OpenAndReadFileContent(abs string) ([]byte, error)
	Open(path string) (io.ReadCloser, error)
	IsDir(path string) (bool, error)
	Walk(ctx context.Context,
		root string,
		isValid func(info os.FileInfo) bool,
		handlePath func(ctx context.Context, path string) error) *errgroup.Group
}
