package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sbreitf1/fs/path"

	"github.com/sbreitf1/errors"
	"github.com/stretchr/testify/assert"
)

//TODO move to file system tests

func TestExistsDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "fs-test-")
	if err != nil {
		panic(err)
	}

	exists, err := Exists(tmpDir)
	assert.NoError(t, err)
	assert.True(t, exists)

	os.RemoveAll(tmpDir)
	exists, err = Exists(tmpDir)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestExistsFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "fs-test-")
	if err != nil {
		panic(err)
	}
	tmpFileName := tmpFile.Name()
	tmpFile.Close()

	exists, err := Exists(tmpFileName)
	assert.NoError(t, err)
	assert.True(t, exists)

	os.Remove(tmpFileName)
	exists, err = Exists(tmpFileName)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCopyFile(t *testing.T) {
	assert.NoError(t, WithTempDir("fs-test-", func(tmpDir string) errors.Error {
		expectedData := []byte("this is a test")

		oldFile := filepath.Join(tmpDir, "test.txt")
		ioutil.WriteFile(oldFile, expectedData, os.ModePerm)
		newFile := filepath.Join(tmpDir, "other.txt")

		exists, err := Exists(newFile)
		assert.NoError(t, err)
		assert.False(t, exists)

		assert.NoError(t, CopyFile(oldFile, newFile))
		exists, err = Exists(newFile)
		assert.NoError(t, err)
		assert.True(t, exists)

		data, readErr := ioutil.ReadFile(newFile)
		assert.NoError(t, readErr)
		assert.Equal(t, expectedData, data)

		return nil
	}))
}

func TestReadLines(t *testing.T) {
	assert.NoError(t, WithTempDir("fs-test-", func(tmpDir string) errors.Error {
		ioutil.WriteFile(path.Join(tmpDir, "lines-test.txt"), []byte("this\nis\r\n\na new line with spaces\r"), os.ModePerm)

		lines, err := ReadLines(path.Join(tmpDir, "lines-test.txt"))
		assert.NoError(t, err)
		assert.Equal(t, []string{"this", "is", "", "a new line with spaces", ""}, lines)

		return nil
	}))
}

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
