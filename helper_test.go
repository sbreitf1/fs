package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.NoError(t, WithTempDir("fs-test-", func(tmpDir string) error {
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
	lines, err := ReadLines("test/lines-test.txt")
	assert.NoError(t, err)
	assert.Equal(t, []string{"this", "is", "", "a new line with spaces", ""}, lines)
}
