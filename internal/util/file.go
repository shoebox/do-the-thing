package util

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func NewFileService() IoUtilFileService {
	return IoUtilFileService{}
}

type IoUtilFileService struct {
}

func (f IoUtilFileService) Walk(ctx context.Context,
	root string,
	isValid func(info os.FileInfo) bool,
	handlePath func(ctx context.Context, path string) error) *errgroup.Group {

	// Create an error group for the context
	g, ctx := errgroup.WithContext(ctx)

	// Create a channel to host the paths
	paths := make(chan string)

	// Iterate on all items contained in the target root path
	g.Go(func() error {
		// closing the channel on defer
		defer close(paths)

		// We walk the root folder files
		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			// Do the file is a certificate file (at first glance at least)
			if err != nil || !isValid(info) {
				return nil
			}

			// Select action for result
			select {
			case paths <- path: // Populate the channel with result
			case <-ctx.Done(): // In case of context cancelation
				return ctx.Err()
			}
			return nil
		})
	})

	// We iterate on all paths contained in the channel
	for path := range paths {
		// And launch a goroutine against it
		g.Go(func() error {
			return handlePath(ctx, path)
		})
	}

	return g
}

func (f IoUtilFileService) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (f IoUtilFileService) OpenAndReadFileContent(abs string) ([]byte, error) {
	// Open the file content
	file, err := os.Open(abs)
	if err != nil {
		return nil, err
	}

	// Read the file content
	return ioutil.ReadAll(file)
}

func (f IoUtilFileService) IsDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return stat.IsDir(), nil
}

// TempFileName Generate a temporary file path
func TempFileName(prefix, suffix string) (string, error) {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix), nil
}
