package interop

import (
	"testing"

	"github.com/sbreitf1/errors"
	"github.com/sbreitf1/fs"
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
		errors.AssertNil(t, Copy(fs1, "/test.txt", fs2, "/out.txt"))
		assertIsFile(t, fs1, "/test.txt")
		assertFileContent(t, fs2, "/out.txt", "foo bar")
	})

	t.Run("TestCopyDir", func(t *testing.T) {
		errors.AssertNil(t, fs1.CreateDirectory("/foo"))
		errors.AssertNil(t, fs1.CreateDirectory("/foo/bar"))
		errors.AssertNil(t, fs1.CreateDirectory("/foo/test"))
		errors.AssertNil(t, fs1.CreateDirectory("/foo/bar/hello"))
		errors.AssertNil(t, fs1.WriteString("/foo/test.txt", "foo1"))
		errors.AssertNil(t, fs1.WriteString("/foo/bar/hello/blub.txt", "bar2"))
		errors.AssertNil(t, Copy(fs1, "/foo", fs2, "/nice"))
		assertIsDir(t, fs1, "/foo")
		assertFileContent(t, fs2, "/nice/test.txt", "foo1")
		assertFileContent(t, fs2, "/nice/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/nice/test")
	})
}
