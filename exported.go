package fs

import (
	"github.com/sbreitf1/errors"
)

var (
	// DefaultFileSystem denots the file system used for all default accessors.
	DefaultFileSystem *FileSystem
)

func init() {
	DefaultFileSystem = New()
}

// Exists returns true, if the given path is a file or directory.
func Exists(path string) (bool, errors.Error) {
	return DefaultFileSystem.Exists(path)
}

// IsFile returns true, if the given path is a file.
func IsFile(path string) (bool, errors.Error) {
	return DefaultFileSystem.IsFile(path)
}

// IsDir returns true, if the given path is a directory.
func IsDir(path string) (bool, errors.Error) {
	return DefaultFileSystem.IsDir(path)
}

// ReadLines returns all lines separated by "\n", "\r" or "\r\n" from a file.
func ReadLines(path string) ([]string, errors.Error) {
	return DefaultFileSystem.ReadLines(path)
}

// WriteLines writes all lines separated by the default line separator to a file.
func WriteLines(path string, lines []string) errors.Error {
	return DefaultFileSystem.WriteLines(path, lines)
}

// Copy clone a file or directory to the target. If the target already exists, it must be the same element type (file or directory) to be overwritten.
func Copy(src, dst string) errors.Error {
	panic("Copy not implemented yet")
}

// CopyFile clones a file and overwrites the existing one.
func CopyFile(src, dst string) errors.Error {
	return DefaultFileSystem.CopyFile(src, dst)
}

// CopyDir recursively clones a directory overwriting all existing files.
func CopyDir(src, dst string) errors.Error {
	panic("CopyDir not implemented yet")
}

// WithTempFile creates a temporary file and deletes it when f returns.
func WithTempFile(pattern string, f func(tmpFile string) errors.Error) errors.Error {
	return DefaultFileSystem.WithTempFile(pattern, f)
}

// WithTempDir creates a temporary directory and deletes it when f returns.
func WithTempDir(prefix string, f func(tmpDir string) errors.Error) errors.Error {
	return DefaultFileSystem.WithTempDir(prefix, f)
}
