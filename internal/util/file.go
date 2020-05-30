package util

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileService interface {
	OpenAndReadFileContent(abs string) ([]byte, error)
	IsDir(path string) (bool, error)
}

type IoUtilFileService struct {
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
