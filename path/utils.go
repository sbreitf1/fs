package path

import (
	"path/filepath"

	"github.com/sbreitf1/errors"
)

const (
	// DefaultPathDelimiter denotes the character used to separate directory and file names in paths.
	DefaultPathDelimiter = "/"
)

var (
	// ErrInvalidPath occurs when using malformed paths.
	ErrInvalidPath = errors.New("Malformed path")
)

// Join merges multiple path parts using the DefaultPathDelimiter.
func Join(parts ...string) string {
	return filepath.Join(parts...)
}

// Base returns only the last part of a path.
func Base(path string) string {
	return filepath.Base(path)
}

// BaseNoExt returns only the last part of a path excluding the file extension as returned by Ext.
func BaseNoExt(path string) string {
	return Base(NoExt(path))
}

// Dir returns the parent directory of a path.
func Dir(path string) string {
	return filepath.Dir(path)
}

// Abs retrieves the full path to a relative path.
func Abs(path string) (string, errors.Error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", ErrInvalidPath.Make().Cause(err)
	}
	return path, nil
}

// Ext returns the file extensions including the dot character.
func Ext(path string) string {
	return filepath.Ext(path)
}

// NoExt removes the file extension as returned by Ext.
func NoExt(path string) string {
	ext := Ext(path)
	return path[:len(path)-len(ext)]
}
