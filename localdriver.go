package fs

import (
	"io/ioutil"
	"os"

	"github.com/sbreitf1/fs/path"

	"github.com/sbreitf1/errors"
)

// LocalDriver allows access to the file system of the host machine.
type LocalDriver struct {
	Root string
}

func (d *LocalDriver) root(p string) (string, errors.Error) {
	if len(d.Root) == 0 {
		return p, nil
	}
	if !path.IsAbs(p) {
		return "", path.Err.Msg("Relative paths are not allowed on rooted local file systems").Make()
	}
	return path.AbsRoot(d.Root, p)
}

// Exists returns true, if the given path is a file or directory.
func (d *LocalDriver) Exists(path string) (bool, errors.Error) {
	rootedPath, err := d.root(path)
	if err != nil {
		return false, err
	}

	_, statErr := os.Stat(rootedPath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

// IsFile returns true, if the given path is a file.
func (d *LocalDriver) IsFile(path string) (bool, errors.Error) {
	rootedPath, err := d.root(path)
	if err != nil {
		return false, err
	}

	fi, statErr := os.Stat(rootedPath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return false, nil
		}
		return false, nil
	}
	return !fi.IsDir(), nil
}

// IsDir returns true, if the given path is a directory.
func (d *LocalDriver) IsDir(path string) (bool, errors.Error) {
	rootedPath, err := d.root(path)
	if err != nil {
		return false, err
	}

	fi, statErr := os.Stat(rootedPath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return false, nil
		}
		return false, nil
	}
	return fi.IsDir(), nil
}

// ReadDir returns all files and directories contained in a directory.
func (d *LocalDriver) ReadDir(path string) ([]FileInfo, errors.Error) {
	rootedPath, err := d.root(path)
	if err != nil {
		return nil, err
	}

	items, readErr := ioutil.ReadDir(rootedPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, ErrDirectoryNotExists.Msg("Directory %q not found", path).Make()
		}
		return nil, Err.Msg("Failed to list directory content").Make().Cause(readErr)
	}

	result := make([]FileInfo, len(items))
	for i := range items {
		result[i] = items[i]
	}
	return result, nil
}

// OpenFile opens a file instance and returns the handle.
func (d *LocalDriver) OpenFile(path string, flags OpenFlags) (File, errors.Error) {
	rootedPath, err := d.root(path)
	if err != nil {
		return nil, err
	}

	f, openErr := os.OpenFile(rootedPath, int(flags), os.ModePerm)
	if openErr != nil {
		if os.IsNotExist(openErr) {
			return nil, ErrFileNotExists.Args(path).Make()
		}
		return nil, Err.Msg("Could not open file").Make().Cause(openErr)
	}
	return f, nil
}

// CreateDirectory creates a new directory and all parent directories if they do not exist.
func (d *LocalDriver) CreateDirectory(path string) errors.Error {
	rootedPath, err := d.root(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(rootedPath, os.ModePerm); err != nil {
		return Err.Msg("Failed to create directory").Make().Cause(err)
	}
	return nil
}

// DeleteFile deletes a file.
func (d *LocalDriver) DeleteFile(path string) errors.Error {
	rootedPath, err := d.root(path)
	if err != nil {
		return err
	}

	if err := errors.Wrap(os.Remove(rootedPath)); err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists.Args(path).Make()
		}
		return Err.Msg("Could not delete file").Make().Cause(err)
	}
	return nil
}

// DeleteDirectory deletes an empty directory. Set recursive to true to also remove directory content.
func (d *LocalDriver) DeleteDirectory(path string, recursive bool) errors.Error {
	rootedPath, err := d.root(path)
	if err != nil {
		return err
	}

	var removeErr error
	if recursive {
		removeErr = os.RemoveAll(rootedPath)
	} else {
		removeErr = os.Remove(rootedPath)
	}

	if removeErr != nil {
		if os.IsNotExist(removeErr) {
			return ErrFileNotExists.Args(path).Make()
		} else if os.IsExist(removeErr) {
			return ErrNotEmpty.Make()
		}
		return Err.Msg("Could not delete directory").Make().Cause(removeErr)
	}
	return nil
}

// MoveFile moves a file to a new location.
func (d *LocalDriver) MoveFile(src, dst string) errors.Error {
	rootedSrc, err := d.root(src)
	if err != nil {
		return err
	}
	rootedDst, err := d.root(dst)
	if err != nil {
		return err
	}

	if err := os.Rename(rootedSrc, rootedDst); err != nil {
		if os.IsNotExist(err) {
			//TODO check src file or dst dir does not exist
			return ErrFileNotExists.Args(src).Make()
		}
		return Err.Msg("Could not move file").Make().Cause(err)
	}
	return nil
}

// MoveDir moves a directory to a new location.
func (d *LocalDriver) MoveDir(src, dst string) errors.Error {
	rootedSrc, err := d.root(src)
	if err != nil {
		return err
	}
	rootedDst, err := d.root(dst)
	if err != nil {
		return err
	}

	if err := os.Rename(rootedSrc, rootedDst); err != nil {
		if os.IsNotExist(err) {
			//TODO check src file or dst dir does not exist
			return ErrFileNotExists.Args(src).Make()
		}
		return Err.Msg("Could not move directory").Make().Cause(err)
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
		return "", Err.Msg("Failed to create temporary file").Make().Cause(err)
	}
	defer tmpFile.Close()
	return tmpFile.Name(), nil
}

// GetTempDir returns the path to an empty temporary dir.
func (d *LocalDriver) GetTempDir(prefix string) (string, errors.Error) {
	if len(d.Root) > 0 {
		return "", ErrNotSupported.Msg("Cannot create temporary directories on rooted local file systems").Make()
	}

	tmpDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		return "", Err.Msg("Failed to create temporary directory").Make().Cause(err)
	}
	return tmpDir, nil
}
