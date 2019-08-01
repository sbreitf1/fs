package fs

import (
	"testing"

	"github.com/sbreitf1/errors"
	"github.com/stretchr/testify/assert"
)

//TODO move to file system tests
func TestWithTempDir(t *testing.T) {
	var dir string
	assert.NoError(t, WithTempDir("fs-test-", func(tmpDir string) errors.Error {
		dir = tmpDir

		exists, err := IsDir(dir)
		if err != nil {
			panic(err)
		}
		assert.True(t, exists, "Directory should have been created")

		return nil
	}))
	exists, err := IsDir(dir)
	if err != nil {
		panic(err)
	}
	assert.False(t, exists, "Directory should have been deleted")
}

func TestWithTempFile(t *testing.T) {
	var file string
	assert.NoError(t, WithTempFile("fs-test-", func(tmpFile string) errors.Error {
		file = tmpFile

		exists, err := IsFile(file)
		if err != nil {
			panic(err)
		}
		assert.True(t, exists, "File should have been created")

		return nil
	}))
	exists, err := IsFile(file)
	if err != nil {
		panic(err)
	}
	assert.False(t, exists, "File should have been deleted")
}
