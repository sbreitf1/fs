package interop

import (
	"testing"

	"github.com/sbreitf1/fs"

	"github.com/sbreitf1/errors"
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
		errors.AssertNil(t, MoveFile(fs1, "/test.txt", fs2, "/out.txt"))
		assertNotExists(t, fs1, "/test.txt")
		assertFileContent(t, fs2, "/out.txt", "foo bar")
	})

	t.Run("TestSelectMoveFile", func(t *testing.T) {
		errors.AssertNil(t, fs1.WriteString("/test.txt", "foo bar"))
		errors.AssertNil(t, Move(fs1, "/test.txt", fs2, "/out2.txt"))
		assertNotExists(t, fs1, "/test.txt")
		assertFileContent(t, fs2, "/out2.txt", "foo bar")
	})

	t.Run("TestMoveDir", func(t *testing.T) {
		prepareDir(t, fs1)
		errors.AssertNil(t, MoveDir(fs1, "/foo", fs2, "/nice"))
		assertNotExists(t, fs1, "/foo")
		assertFileContent(t, fs2, "/nice/test.txt", "foo1")
		assertFileContent(t, fs2, "/nice/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/nice/test")
	})

	t.Run("TestSelectMoveDir", func(t *testing.T) {
		prepareDir(t, fs1)
		errors.AssertNil(t, Move(fs1, "/foo", fs2, "/nice"))
		assertNotExists(t, fs1, "/foo")
		assertFileContent(t, fs2, "/nice/test.txt", "foo1")
		assertFileContent(t, fs2, "/nice/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/nice/test")
	})

	t.Run("TestMoveAll", func(t *testing.T) {
		prepareDir(t, fs1)
		errors.AssertNil(t, MoveAll(fs1, "/foo", fs2, "/"))
		assertIsDir(t, fs1, "/foo")
		assertNotExists(t, fs1, "/foo/bar")
		assertNotExists(t, fs1, "/foo/test")
		assertFileContent(t, fs2, "/test.txt", "foo1")
		assertFileContent(t, fs2, "/bar/hello/blub.txt", "bar2")
		assertIsDir(t, fs2, "/bar")
		assertIsDir(t, fs2, "/test")
	})
}
