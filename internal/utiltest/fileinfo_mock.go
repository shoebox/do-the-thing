package utiltest

import "os"

// MockFileInfo is a mock of the os.FileInfo class to be able to test different configuration
type MockFileInfo struct {
	os.FileInfo             // Embed this so we only need to add methods used by testable functions
	fileMode    os.FileMode // mode
	isDir       bool        // Do the file is a directory
	name        string      // file base name
}

func NewMockFileInfo(fm os.FileMode, dir bool, n string) MockFileInfo {
	return MockFileInfo{
		fileMode: fm,
		name:     n,
		isDir:    dir,
	}
}

// Name will return the configured value for the mock of the name field
func (m MockFileInfo) Name() string { return m.name }

// Mode will return the configure value for the mock of the fileMode field
func (m MockFileInfo) Mode() os.FileMode { return m.fileMode }

// IsDir will return the configure value for the mock of the isDir field
func (m MockFileInfo) IsDir() bool { return m.isDir }
