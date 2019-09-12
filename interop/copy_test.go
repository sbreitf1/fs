package interop

import (
	"testing"

	"github.com/sbreitf1/fs"

	"github.com/sbreitf1/errors"
)

func TestCopy(t *testing.T) {
	fs.WithTempDir("fs-test-", func(tmpDir1 string) errors.Error {
		return fs.WithTempDir("fs-test-", func(tmpDir2 string) errors.Error {
			fs1 := fs.NewWithDriver(&fs.LocalDriver{Root: tmpDir1})
			fs2 := fs.NewWithDriver(&fs.LocalDriver{Root: tmpDir2})
			testCopy(t, fs1, fs2)
			return nil
		})
	})
}

func testCopy(t *testing.T, fs1, fs2 *fs.FileSystem) {
	t.Run("TestCopyFile", func(t *testing.T) {
		errors.AssertNil(t, fs1.WriteString("/test.txt", "foo bar"))
		errors.AssertNil(t, CopyFile(fs1, "/test.txt", fs2, "/out.txt"))
		assertIsFile(t, fs1, "/test.txt")
		assertFileContent(t, fs2, "/out.txt", "foo bar")
	})

	t.Run("TestSelectCopyFile", func(t *testing.T) {
		errors.AssertNil(t, fs1.WriteString("/test.txt", "foo bar"))
		errors.AssertNil(t, Copy(fs1, "/test.txt", fs2, "/out2.txt"))
		assertIsFile(t, fs1, "/test.txt")
		assertFileContent(t, fs2, "/out2.txt", "foo bar")
	})

	t.Run("TestCopyDir", func(t *testing.T) {
		prepareDir(t, fs1)
		errors.AssertNil(t, CopyDir(fs1, "/foo", fs2, "/nice"))
		assertIsDir(t, fs1, "/foo")
		assertFileContent(t, fs2, "/nice/test.txt", "foo1")
		assertFileContent(t, fs2, "/nice/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/nice/test")
	})

	t.Run("TestSelectCopyDir", func(t *testing.T) {
		errors.AssertNil(t, Copy(fs1, "/foo", fs2, "/nice2"))
		assertIsDir(t, fs1, "/foo")
		assertFileContent(t, fs2, "/nice2/test.txt", "foo1")
		assertFileContent(t, fs2, "/nice2/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/nice2/test")
	})

	t.Run("TestCopyAll", func(t *testing.T) {
		errors.AssertNil(t, CopyAll(fs1, "/foo", fs2, "/"))
		assertIsDir(t, fs1, "/foo")
		assertFileContent(t, fs2, "/test.txt", "foo1")
		assertFileContent(t, fs2, "/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/bar")
		assertIsDir(t, fs2, "/test")
	})
}
