package path

import (
	"path/filepath"
)

const (
	// DefaultPathDelimiter denotes the character used to separate directory and file names in paths.
	DefaultPathDelimiter = "/"
)

// Join merges multiple path parts using the DefaultPathDelimiter.
func Join(parts ...string) string {
	return filepath.Join(parts...)
}
