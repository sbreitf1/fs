package fs

import (
	"io"
	"os"
)

// Exists returns true, if the given path is a file or directory.
func Exists(path string) (bool, error) {
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
func IsFile(path string) (bool, error) {
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
func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, nil
	}
	return fi.IsDir(), nil
}

// Copy clone a file or directory to the target. If the target already exists, it must be the same element type (file or directory) to be overwritten.
func Copy(src, dst string) error {
	panic("Copy not implemented yet")
}

// CopyFile clones a file and overwrites the existing one.
func CopyFile(src, dst string) error {
	//TODO handle overwrite?

	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	return err
}

// CopyDir recursively clones a directory overwriting all existing files.
func CopyDir(src, dst string) error {
	panic("CopyDir not implemented yet")
}
