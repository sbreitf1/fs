package interop

import (
	"testing"

	"github.com/sbreitf1/errors"
	"github.com/sbreitf1/fs"
	"github.com/stretchr/testify/assert"
)

/* ############################################### */
/* ###               Test Helper               ### */
/* ############################################### */

func assertNotExists(t *testing.T, fs *fs.FileSystem, path string) bool {
	exists, err := fs.Exists(path)
	if errors.AssertNil(t, err, "Error while checking for %q", path) {
		return assert.False(t, exists, "Expected %q to not exist", path)
	}
	return false
}

func assertIsFile(t *testing.T, fs *fs.FileSystem, path string) bool {
	isFile, err := fs.IsFile(path)
	if errors.AssertNil(t, err, "Error while checking for file %q", path) {
		return assert.True(t, isFile, "Expected file %q does not exist", path)
	}
	return false
}

func assertIsDir(t *testing.T, fs *fs.FileSystem, path string) bool {
	isDir, err := fs.IsDir(path)
	if errors.AssertNil(t, err, "Error while checking for dir %q", path) {
		return assert.True(t, isDir, "Expected directory %q does not exist", path)
	}
	return false
}

func assertFileContent(t *testing.T, fs *fs.FileSystem, path, expectedContent string) bool {
	data, err := fs.ReadString(path)
	if errors.AssertNil(t, err, "Error while accessing fiile %q", path) {
		return assert.Equal(t, expectedContent, data, "Unexpected file content of %q", path)
	}
	return false
}
