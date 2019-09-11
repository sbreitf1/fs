package interop

import (
	"testing"

	"github.com/sbreitf1/errors"
	"github.com/sbreitf1/fs"
)

func TestMove(t *testing.T) {
	fs.WithTempDir("fs-test-", func(tmpDir1 string) errors.Error {
		return fs.WithTempDir("fs-test-", func(tmpDir2 string) errors.Error {
			fs1 := fs.NewWithDriver(&fs.LocalDriver{Root: tmpDir1})
			fs2 := fs.NewWithDriver(&fs.LocalDriver{Root: tmpDir2})
			testMove(t, fs1, fs2)
			return nil
		})
	})
}

func testMove(t *testing.T, fs1, fs2 *fs.FileSystem) {
	t.Run("TestMoveFile", func(t *testing.T) {
		errors.AssertNil(t, fs1.WriteString("/test.txt", "foo bar"))
		errors.AssertNil(t, Move(fs1, "/test.txt", fs2, "/out.txt"))
		assertNotExists(t, fs1, "/test.txt")
		assertFileContent(t, fs2, "/out.txt", "foo bar")
	})

	t.Run("TestMoveDir", func(t *testing.T) {
		errors.AssertNil(t, fs1.CreateDirectory("/foo"))
		errors.AssertNil(t, fs1.CreateDirectory("/foo/bar"))
		errors.AssertNil(t, fs1.CreateDirectory("/foo/test"))
		errors.AssertNil(t, fs1.CreateDirectory("/foo/bar/hello"))
		errors.AssertNil(t, fs1.WriteString("/foo/test.txt", "foo1"))
		errors.AssertNil(t, fs1.WriteString("/foo/bar/hello/blub.txt", "bar2"))
		errors.AssertNil(t, Move(fs1, "/foo", fs2, "/nice"))
		assertNotExists(t, fs1, "/foo")
		assertFileContent(t, fs2, "/nice/test.txt", "foo1")
		assertFileContent(t, fs2, "/nice/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/nice/test")
	})
}
