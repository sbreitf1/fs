package fs

import (
	"os"

	"github.com/sbreitf1/errors"
)

// LocalDriver allows access to the file system of the host machine.
type LocalDriver struct {
	//TODO respect root dir in methods
	Root string
}

// Exists returns true, if the given path is a file or directory.
func (d *LocalDriver) Exists(path string) (bool, errors.Error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

// IsFile returns true, if the given path is a file.
func (d *LocalDriver) IsFile(path string) (bool, errors.Error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, nil
	}
	return !fi.IsDir(), nil
}

// IsDir returns true, if the given path is a directory.
func (d *LocalDriver) IsDir(path string) (bool, errors.Error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, nil
	}
	return fi.IsDir(), nil
}

func (d *LocalDriver) Open(path string) (File, errors.Error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExists.Args(path).Make()
		}
		return nil, errors.Wrap(err).Expand("Could not open file")
	}
	return f, nil
}

func (d *LocalDriver) GetTempFile(prefix string) (string, errors.Error) {
	return "", ErrNotSupported.Args("GetTempFile").Make()
}

func (d *LocalDriver) Create(path string) (File, errors.Error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err).Expand("Could not create file")
	}
	return f, nil
}
