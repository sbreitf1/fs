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
	var fs *Util
	assert.NotPanics(t, func() { fs = New() })
	assert.True(t, fs.CanRead(), "CanRead() returns false")
	assert.True(t, fs.CanWrite(), "CanRead() returns false")
	assert.True(t, fs.CanReadWrite(), "CanReadWrite() returns false")
	assert.True(t, fs.CanTemp(), "CanTemp() returns false")
	assert.True(t, fs.CanAll(), "CanAll() returns false")
}

func TestNewUtilInvalid(t *testing.T) {
	assert.Panics(t, func() { NewUtil("not a file system driver") })
}

func TestReadString(t *testing.T) {
	WithTempDir("fs-test-", func(tmpFile string) error {
		fs := New()
		path := filepath.Join(tmpFile, "test.txt")
		if err := ioutil.WriteFile(path, []byte("a new cool file content"), os.ModePerm); err != nil {
			panic(err)
		}
		data, err := fs.ReadString(path)
		errors.AssertNil(t, err)
		assert.Equal(t, "a new cool file content", data)
		return nil
	})
}

func TestWriteLines(t *testing.T) {
	WithTempDir("fs-test-", func(tmpFile string) error {
		fs := New()
		path := filepath.Join(tmpFile, "test.txt")
		errors.AssertNil(t, fs.WriteLines(path, []string{"foo", "bar", "", "yeah!", ""}))
		assert.FileExists(t, path)
		data, err := fs.ReadString(path)
		errors.AssertNil(t, err)
		assert.Equal(t, "foo\nbar\n\nyeah!\n", data)
		return nil
	})
}
