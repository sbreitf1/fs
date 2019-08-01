package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sbreitf1/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var fs *FileSystem
	assert.NotPanics(t, func() { fs = New() })
	assert.True(t, fs.CanRead(), "CanRead() returns false")
	assert.True(t, fs.CanWrite(), "CanRead() returns false")
	assert.True(t, fs.CanReadWrite(), "CanReadWrite() returns false")
	assert.True(t, fs.CanTemp(), "CanTemp() returns false")
	assert.True(t, fs.CanAll(), "CanAll() returns false")
}

func TestNewUtilInvalid(t *testing.T) {
	assert.Panics(t, func() { NewWithDriver("not a file system driver") })
}

func TestFileSystemFunctionality(t *testing.T) {
	fs := New()
	errors.AssertNil(t, WithTempDir("fs-test-", func(tmpDir string) errors.Error {
		testFS(t, fs, tmpDir)
		return nil
	}))
}

func testFS(t *testing.T, fs *FileSystem, tmpDir string) {
	t.Run("TestReadString", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.txt")
		if err := ioutil.WriteFile(path, []byte("a new cool file content"), os.ModePerm); err != nil {
			panic(err)
		}
		data, err := fs.ReadString(path)
		errors.AssertNil(t, err)
		assert.Equal(t, "a new cool file content", data)
	})

	t.Run("TestWriteLines", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.txt")
		errors.AssertNil(t, fs.WriteLines(path, []string{"foo", "bar", "", "yeah!", ""}))
		assert.FileExists(t, path)
		data, err := fs.ReadString(path)
		errors.AssertNil(t, err)
		assert.Equal(t, "foo\nbar\n\nyeah!\n", data)
	})
}
