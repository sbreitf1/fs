package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTempDir(t *testing.T) {
	var dir string
	assert.NoError(t, WithTempDir("fs-test-", func(tmpDir string) error {
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
	assert.NoError(t, WithTempFile("fs-test-", func(tmpFile string) error {
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
