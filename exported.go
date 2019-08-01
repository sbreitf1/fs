package fs

import (
	"github.com/sbreitf1/errors"
)

var (
	// DefaultFileSystem denots the file system used for all default accessors.
	DefaultFileSystem *Util
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
