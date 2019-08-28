package fs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sbreitf1/fs/path"

	"github.com/sbreitf1/errors"
	"github.com/stretchr/testify/assert"
)

func TestOpenFlags(t *testing.T) {
	assert.True(t, OpenReadOnly.IsRead())
	assert.False(t, OpenReadOnly.IsWrite())
	assert.False(t, OpenWriteOnly.IsRead())
	assert.True(t, OpenWriteOnly.IsWrite())
	assert.True(t, OpenReadWrite.IsRead())
	assert.True(t, OpenReadWrite.IsWrite())
	assert.False(t, OpenReadOnly.Append().Create().Exclusive().Sync().Truncate().IsWrite())
	assert.Equal(t, OpenWriteOnly, OpenWriteOnly.Exclusive().Access())
}

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
	assert.Panics(t, func() { NewWithDriver(nil) })
	assert.Panics(t, func() { NewWithDriver("not a file system driver") })
}

func TestFileSystemCommon(t *testing.T) {
	fs := New()
	errors.AssertNil(t, WithTempDir("fs-test-", func(tmpDir string) errors.Error {
		testFS(t, fs, tmpDir)
		return nil
	}))
}

func testFS(t *testing.T, fs *FileSystem, dir string) {
	t.Run("TestReadString", func(t *testing.T) {
		path := path.Join(dir, "test.txt")
		if err := ioutil.WriteFile(path, []byte("a new cool file content"), os.ModePerm); err != nil {
			panic(err)
		}
		assertFileContent(t, fs, path, "a new cool file content")
	})

	t.Run("TestWriteLines", func(t *testing.T) {
		path := path.Join(dir, "test.txt")
		errors.AssertNil(t, fs.WriteLines(path, []string{"foo", "bar", "", "yeah!", ""}))
		assertFileContent(t, fs, path, "foo\nbar\n\nyeah!\n")
	})

	t.Run("TestCreateDirectory", func(t *testing.T) {
		path := path.Join(dir, "testdir/subdir")
		assertNotExists(t, fs, path)
		fs.CreateDirectory(path)
		assertIsDir(t, fs, path)
	})

	t.Run("TestCopyFile", func(t *testing.T) {
		src := path.Join(dir, "test.txt")
		dst := path.Join(dir, "testdir/subdir/foobar.txt")
		errors.AssertNil(t, fs.CopyFile(src, dst))
		assertFileContent(t, fs, dst, "foo\nbar\n\nyeah!\n")
	})

	t.Run("TestCopyDir", func(t *testing.T) {
		src := path.Join(dir, "testdir")
		dst := path.Join(dir, "justanotherdir")
		fs.CreateDirectory(path.Join(src, "subdir/foobar1337"))
		errors.AssertNil(t, fs.CopyDir(src, dst))
		// new file has correct content
		assertFileContent(t, fs, path.Join(dst, "subdir/foobar.txt"), "foo\nbar\n\nyeah!\n")
		// empty dir is copied aswell
		assertIsDir(t, fs, path.Join(dst, "subdir/foobar1337"))
		// old file still exists
		assertIsFile(t, fs, path.Join(src, "subdir/foobar.txt"))
	})

	t.Run("TestCopyAll", func(t *testing.T) {
		src := path.Join(dir, "justanotherdir/subdir")
		errors.AssertNil(t, fs.CopyAll(src, dir))
		// new file has correct content
		assertFileContent(t, fs, path.Join(dir, "foobar.txt"), "foo\nbar\n\nyeah!\n")
		// old file still exists
		assertIsFile(t, fs, path.Join(src, "foobar.txt"))
	})

	t.Run("TestOpenWrite", func(t *testing.T) {
		path := path.Join(dir, "openwritetest.txt")
		errors.AssertNil(t, fs.WriteString(path, "foo bar cool test data content"))

		f, err := fs.OpenFile(path, OpenReadWrite)
		errors.AssertNil(t, err)
		f.Write([]byte("short stuff"))
		f.Close()

		assertFileContent(t, fs, path, "short stuffl test data content")
	})

	t.Run("TestTruncate", func(t *testing.T) {
		path := path.Join(dir, "trunctest.txt")
		errors.AssertNil(t, fs.WriteString(path, "foo bar cool test data content"))

		f, err := fs.OpenFile(path, OpenReadWrite.Truncate())
		errors.AssertNil(t, err)
		f.Write([]byte("short stuff"))
		f.Close()

		assertFileContent(t, fs, path, "short stuff")
	})

	t.Run("TestAppend", func(t *testing.T) {
		path := path.Join(dir, "appendtest.txt")
		errors.AssertNil(t, fs.WriteString(path, "foo bar"))

		f, err := fs.OpenFile(path, OpenReadWrite.Append())
		errors.AssertNil(t, err)
		f.Write([]byte(" - short stuff"))
		f.Close()

		assertFileContent(t, fs, path, "foo bar - short stuff")
	})
}

/* ############################################### */
/* ###               Test Heper                ### */
/* ############################################### */

func assertNotExists(t *testing.T, fs *FileSystem, path string) bool {
	exists, err := fs.Exists(path)
	if errors.AssertNil(t, err, "Error while checking for %q", path) {
		return assert.False(t, exists, "Expected %q to not exist", path)
	}
	return false
}

func assertIsFile(t *testing.T, fs *FileSystem, path string) bool {
	isFile, err := fs.IsFile(path)
	if errors.AssertNil(t, err, "Error while checking for file %q", path) {
		return assert.True(t, isFile, "Expected file %q does not exist", path)
	}
	return false
}

func assertIsDir(t *testing.T, fs *FileSystem, path string) bool {
	isDir, err := fs.IsDir(path)
	if errors.AssertNil(t, err, "Error while checking for dir %q", path) {
		return assert.True(t, isDir, "Expected directory %q does not exist", path)
	}
	return false
}

func assertFileContent(t *testing.T, fs *FileSystem, path, expectedContent string) bool {
	data, err := fs.ReadString(path)
	if errors.AssertNil(t, err, "Error while accessing fiile %q", path) {
		return assert.Equal(t, expectedContent, data, "Unexpected file content of %q", path)
	}
	return false
}
