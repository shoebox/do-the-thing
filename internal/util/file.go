package util

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func NewFileService() IoUtilFileService {
	return IoUtilFileService{}
}

type IoUtilFileService struct {
}

func (f IoUtilFileService) Walk(
	ctx context.Context,
	root string,
	isValid func(info os.FileInfo) bool,
	paths chan string,
	wg *sync.WaitGroup,
) error {

	defer close(paths)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
		if err != nil {
			return err
		}

		if !isValid(info) {
			return nil
		}

		paths <- path

		return nil
	})

	return err
}

func (f IoUtilFileService) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// TODO: Return a reader rather
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
