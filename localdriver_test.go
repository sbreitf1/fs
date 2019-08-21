package fs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sbreitf1/fs/path"

	"github.com/sbreitf1/errors"
	"github.com/stretchr/testify/assert"
)

func TestLocalDriverCommon(t *testing.T) {
	t.Run("TestLocalDriver", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "fs-test-")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(tmpDir)
		testLocalDriver(t, &LocalDriver{}, "", tmpDir)
	})

	t.Run("TestRootedLocalDriver", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "fs-test-")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(tmpDir)
		testLocalDriver(t, &LocalDriver{Root: tmpDir}, tmpDir, "")
	})
}

func testLocalDriver(t *testing.T, driver *LocalDriver, rootDir, workingDir string) {
	t.Run("ReadDirEmpty", func(t *testing.T) {
		files, err := driver.ReadDir(path.Join(workingDir, "/"))
		errors.AssertNil(t, err)
		assert.Equal(t, 0, len(files))
	})

	t.Run("ReadDirNonExistent", func(t *testing.T) {
		_, err := driver.ReadDir(path.Join(workingDir, "/nonexistingpath"))
		errors.Assert(t, ErrDirectoryNotExists, err)
	})

	t.Run("IsFile", func(t *testing.T) {
		if err := ioutil.WriteFile(path.Join(rootDir, workingDir, "/test.txt"), []byte("test data"), os.ModePerm); err != nil {
			panic(err)
		}

		isFile, err := driver.IsFile(path.Join(workingDir, "/test.txt"))
		errors.AssertNil(t, err)
		assert.True(t, isFile)
	})

	t.Run("OpenFile", func(t *testing.T) {
		f, err := driver.OpenFile(path.Join(workingDir, "/test.txt"), OpenReadOnly)
		defer f.Close()
		errors.AssertNil(t, err)

		data, readErr := ioutil.ReadAll(f)
		errors.AssertNil(t, readErr)
		assert.Equal(t, "test data", string(data))
	})

	t.Run("ReadDirSingleFile", func(t *testing.T) {
		files, err := driver.ReadDir(path.Join(workingDir, "/"))
		errors.AssertNil(t, err)
		assert.Equal(t, 1, len(files))
		assert.Equal(t, "test.txt", files[0].Name())
		assert.False(t, files[0].IsDir())
	})

	t.Run("CreateDir", func(t *testing.T) {
		errors.AssertNil(t, driver.CreateDirectory(path.Join(workingDir, "/newdir/and/subdir")))
		assert.DirExists(t, path.Join(rootDir, workingDir, "/newdir"))
		assert.DirExists(t, path.Join(rootDir, workingDir, "/newdir/and"))
		assert.DirExists(t, path.Join(rootDir, workingDir, "/newdir/and/subdir"))
	})

	t.Run("CreateFile", func(t *testing.T) {
		f, err := driver.OpenFile(path.Join(workingDir, "/newdir/and/subdir/testfile.txt"), OpenReadWrite.Create().Truncate())
		errors.AssertNil(t, err)

		f.Write([]byte("some test data"))
		f.Close()

		assert.FileExists(t, path.Join(rootDir, workingDir, "/newdir/and/subdir/testfile.txt"))
		data, readErr := ioutil.ReadFile(path.Join(rootDir, workingDir, "/newdir/and/subdir/testfile.txt"))
		errors.AssertNil(t, readErr)
		assert.Equal(t, "some test data", string(data))
	})

	t.Run("MoveFile", func(t *testing.T) {
		driver.MoveFile(path.Join(workingDir, "/newdir/and/subdir/testfile.txt"), path.Join(workingDir, "/newdir/and/testfile.txt"))

		_, err := os.Stat(path.Join(rootDir, workingDir, "/newdir/and/subdir/testfile.txt"))
		assert.True(t, os.IsNotExist(err))

		assert.FileExists(t, path.Join(rootDir, workingDir, "/newdir/and/testfile.txt"))
		data, readErr := ioutil.ReadFile(path.Join(rootDir, workingDir, "/newdir/and/testfile.txt"))
		errors.AssertNil(t, readErr)
		assert.Equal(t, "some test data", string(data))
	})

	t.Run("MoveDir", func(t *testing.T) {
		driver.MoveDir(path.Join(workingDir, "/newdir/and"), path.Join(workingDir, "/foo"))

		_, err := os.Stat(path.Join(rootDir, workingDir, "/newdir/and"))
		assert.True(t, os.IsNotExist(err))

		assert.DirExists(t, path.Join(rootDir, workingDir, "/foo"))
		assert.DirExists(t, path.Join(rootDir, workingDir, "/foo/subdir"))
	})

	t.Run("DeleteFile", func(t *testing.T) {
		errors.AssertNil(t, driver.DeleteFile(path.Join(workingDir, "/foo/testfile.txt")))
		_, err := os.Stat(path.Join(rootDir, workingDir, "/foo/testfile.txt"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("DeleteDir", func(t *testing.T) {
		errors.Assert(t, ErrNotEmpty, driver.DeleteDirectory(path.Join(workingDir, "/foo"), false))
		assert.DirExists(t, path.Join(rootDir, workingDir, "/foo"))

		errors.AssertNil(t, driver.DeleteDirectory(path.Join(workingDir, "/foo"), true))
		_, err := os.Stat(path.Join(rootDir, workingDir, "/foo"))
		assert.True(t, os.IsNotExist(err))

		errors.AssertNil(t, driver.DeleteDirectory(path.Join(workingDir, "/newdir"), false))
		_, err = os.Stat(path.Join(rootDir, workingDir, "/newdir"))
		assert.True(t, os.IsNotExist(err))
	})
}
