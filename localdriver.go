package fs

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/sbreitf1/errors"
	"github.com/sbreitf1/fs/path"
)

// LocalDriver allows access to the file system of the host machine.
type LocalDriver struct {
	Root string
}

func (d *LocalDriver) root(p string) string {
	if len(d.Root) == 0 {
		return p
	}
	// error can only occur for inaccessible working directory -> ignore here
	abs, _ := path.Abs(p)
	return path.Join(d.Root, abs)
}

func (d *LocalDriver) unroot(p string) string {
	if len(d.Root) == 0 {
		return p
	}
	if strings.HasPrefix(p, d.Root) {
		return p[len(d.Root):]
	}
	return p
}

// Exists returns true, if the given path is a file or directory.
func (d *LocalDriver) Exists(path string) (bool, errors.Error) {
	_, err := os.Stat(d.root(path))
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
	fi, err := os.Stat(d.root(path))
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
	fi, err := os.Stat(d.root(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, nil
	}
	return fi.IsDir(), nil
}

// ReadDir returns all files and directories contained in a directory.
func (d *LocalDriver) ReadDir(path string) ([]FileInfo, errors.Error) {
	items, err := ioutil.ReadDir(d.root(path))
	if err != nil {
		return nil, ErrFileSystem.Msg("Failed to list directory content").Make().Cause(err)
	}

	result := make([]FileInfo, len(items))
	for i := range items {
		result[i] = items[i]
	}
	return result, nil
}

// OpenFile opens a file instance and returns the handle.
func (d *LocalDriver) OpenFile(path string) (File, errors.Error) {
	f, err := os.Open(d.root(path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExists.Args(path).Make()
		}
		return nil, ErrFileSystem.Msg("Could not open file").Make().Cause(err)
	}
	return f, nil
}

// CreateFile a new file (or truncate an existing) and return the file instance handle.
func (d *LocalDriver) CreateFile(path string) (File, errors.Error) {
	f, err := os.Create(d.root(path))
	if err != nil {
		return nil, ErrFileSystem.Msg("Could not create file").Make().Cause(err)
	}
	return f, nil
}

// DeleteFile deletes a file.
func (d *LocalDriver) DeleteFile(path string) errors.Error {
	if err := errors.Wrap(os.Remove(d.root(path))); err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists.Args(path).Make()
		}
		return ErrFileSystem.Msg("Could not delete file").Make().Cause(err)
	}
	return nil
}

// DeleteDirectory deletes an empty directory. Set recursive to true to also remove directory content.
func (d *LocalDriver) DeleteDirectory(path string, recursive bool) errors.Error {
	var err errors.Error
	if recursive {
		err = errors.Wrap(os.RemoveAll(d.root(path)))
	} else {
		err = errors.Wrap(os.Remove(d.root(path)))
	}

	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists.Args(path).Make()
		}
		return ErrFileSystem.Msg("Could not delete directory").Make().Cause(err)
	}
	return nil
}

// MoveFile moves a file to a new location.
func (d *LocalDriver) MoveFile(src, dst string) errors.Error {
	if err := os.Rename(d.root(src), d.root(dst)); err != nil {
		if os.IsNotExist(err) {
			//TODO check src file or dst dir does not exist
			return ErrFileNotExists.Args(src).Make()
		}
		return ErrFileSystem.Msg("Could not move file").Make().Cause(err)
	}
	return nil
}

// MoveDir moves a directory to a new location.
func (d *LocalDriver) MoveDir(src, dst string) errors.Error {
	if err := os.Rename(d.root(src), d.root(dst)); err != nil {
		if os.IsNotExist(err) {
			//TODO check src file or dst dir does not exist
			return ErrFileNotExists.Args(src).Make()
		}
		return ErrFileSystem.Msg("Could not move directory").Make().Cause(err)
	}
	return nil
}

// GetTempFile returns the path to an empty temporary file.
func (d *LocalDriver) GetTempFile(pattern string) (string, errors.Error) {
	if len(d.Root) > 0 {
		return "", ErrNotSupported.Msg("Cannot create temporary files on rooted local file systems").Make()
	}

	tmpFile, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", ErrFileSystem.Msg("Failed to create temporary file").Make().Cause(err)
	}
	defer tmpFile.Close()
	return d.unroot(tmpFile.Name()), nil
}

// GetTempDir returns the path to an empty temporary dir.
func (d *LocalDriver) GetTempDir(prefix string) (string, errors.Error) {
	if len(d.Root) > 0 {
		return "", ErrNotSupported.Msg("Cannot create temporary directories on rooted local file systems").Make()
	}

	tmpDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		return "", ErrFileSystem.Msg("Failed to create temporary directory").Make().Cause(err)
	}
	return d.unroot(tmpDir), nil
}
