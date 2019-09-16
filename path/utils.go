package path

import (
	"path/filepath"
	"strings"

	"github.com/sbreitf1/errors"
)

const (
	// DefaultPathDelimiter denotes the character used to separate directory and file names in paths.
	DefaultPathDelimiter = filepath.Separator
)

var (
	// Err occurs when using malformed paths.
	Err = errors.New("Malformed path")
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

// Clean removes all navigation parts (. and ..) and removes empty path parts.
func Clean(path string) string {
	return filepath.Clean(path)
}

// IsAbs returns whether the path is absolute.
func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// Abs retrieves the full path to a relative path.
func Abs(path string) (string, errors.Error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", Err.Make().Cause(err)
	}
	return path, nil
}

// AbsIn returns the absolute path as seen by a given working directory. The working directory is ignored for absolute paths.
func AbsIn(wd, path string) (string, errors.Error) {
	if !IsAbs(wd) {
		return "", Err.Msg("The working directory must be an absolute path").Make()
	}

	if IsAbs(path) {
		return path, nil
	}
	return Clean(Join(wd, path)), nil
}

// AbsRoot returns the absolute path and ensures the result to stay in root.
func AbsRoot(root, path string) (string, errors.Error) {
	if !IsAbs(root) {
		return "", Err.Msg("The root directory must be an absolute path").Make()
	}

	abs := Clean(Join(root, Clean(path)))
	if ok, _ := IsIn(abs, root); ok {
		// full path is inside root directory -> all good
		return abs, nil
	}

	// try to force path to be a full path
	abs = Clean(Join(root, Clean("/"+path)))
	if ok, _ := IsIn(abs, root); ok {
		// full path is inside root directory -> all good
		return abs, nil
	}

	return "", Err.Make()
}

// IsIn returns true when the given path is a (recursive) child of expectedParent. This method can be used for security checks.
func IsIn(path, expectedParent string) (bool, errors.Error) {
	if !IsAbs(path) {
		return false, Err.Msg("path must denote an absolute path").Make()
	}

	if !IsAbs(expectedParent) {
		return false, Err.Msg("expectedParent must denote an absolute path").Make()
	}

	parts := strings.Split(Clean(path), "/")
	expectedParts := strings.Split(Clean(expectedParent), "/")

	if len(parts) < len(expectedParts) {
		// expected parent cannot be parent of path
		return false, nil
	}

	for i := range expectedParts {
		if parts[i] != expectedParts[i] {
			return false, nil
		}
	}

	return true, nil
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
