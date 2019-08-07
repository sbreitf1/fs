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
	t.Run("CreateAndReadFile", func(t *testing.T) {
		if err := ioutil.WriteFile(path.Join(rootDir, workingDir, "/test.txt"), []byte("test data"), os.ModePerm); err != nil {
			panic(err)
		}

		isFile, err := driver.IsFile(path.Join(workingDir, "/test.txt"))
		errors.AssertNil(t, err)
		assert.True(t, isFile)
	})
}
