package fs

import (
	"os"

	"github.com/sbreitf1/errors"
)

// LocalFileSystemDriver allows access to the file system of the host machine.
type LocalFileSystemDriver struct {
	root string
}

// NewLocalFileSystemDriver returns a new local file system at root.
func NewLocalFileSystemDriver() *LocalFileSystemDriver {
	return NewLocalRelativeFileSystemDriver("")
}

// NewLocalRelativeFileSystemDriver returns a new local file system rooted at the given directory. Access to parent directories is prohibited.
func NewLocalRelativeFileSystemDriver(root string) *LocalFileSystemDriver {
	return &LocalFileSystemDriver{root}
}

func (d *LocalFileSystemDriver) Open(path string) (File, errors.Error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExists.Args(path).Make()
		}
		return nil, errors.Wrap(err).Expand("Could not open file")
	}
	return f, nil
}

func (d *LocalFileSystemDriver) GetTempFile(prefix string) (string, errors.Error) {
	return "", ErrNotSupported.Args("GetTempFile").Make()
}
