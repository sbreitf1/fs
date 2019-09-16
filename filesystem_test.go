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
	t.Run("TestStatRoot", func(t *testing.T) {
		fi, err := fs.Stat("/")
		errors.AssertNil(t, err)
		assert.Equal(t, "/", fi.Name())
		assert.True(t, fi.IsDir())
	})

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

	t.Run("TestDeleteFile", func(t *testing.T) {
		path := path.Join(dir, "appendtest.txt")
		errors.AssertNil(t, fs.DeleteFile(path))
		assertNotExists(t, fs, path)
	})

	t.Run("TestDeleteDirFail", func(t *testing.T) {
		p := path.Join(dir, "justanotherdir")
		errors.Assert(t, ErrNotEmpty, fs.DeleteDirectory(p, false))
		assertIsDir(t, fs, p)
	})

	t.Run("TestDeleteDirRecursive", func(t *testing.T) {
		p := path.Join(dir, "justanotherdir")
		errors.AssertNil(t, fs.DeleteDirectory(p, true))
		assertNotExists(t, fs, p)
	})

	t.Run("TestCleanDir", func(t *testing.T) {
		errors.AssertNil(t, fs.CleanDir(dir))
		files, err := fs.ReadDir(dir)
		errors.AssertNil(t, err)
		assert.Len(t, files, 0)
	})

	t.Run("TestWalk", func(t *testing.T) {
		errors.AssertNil(t, fs.CreateDirectory(path.Join(dir, "foo")))
		errors.AssertNil(t, fs.CreateDirectory(path.Join(dir, "foo/bar")))
		errors.AssertNil(t, fs.WriteString(path.Join(dir, "foo/bar/test.txt"), "foo bar"))
		assertWalk(t, fs, path.Join(dir, "foo"), nil, []string{"bar", "test.txt"}, []string{"bar"}, []string{"bar"})
	})

	t.Run("TestWalkVisitRoot", func(t *testing.T) {
		assertWalk(t, fs, path.Join(dir, "foo"), &WalkOptions{VisitRootDir: true}, []string{"foo", "bar", "test.txt"}, []string{"bar"}, []string{"bar"})
	})

	t.Run("TestWalkRootCallback", func(t *testing.T) {
		assertWalk(t, fs, path.Join(dir, "foo"), &WalkOptions{EnterLeaveCallbacksForRoot: true}, []string{"bar", "test.txt"}, []string{"foo", "bar"}, []string{"bar", "foo"})
	})

	t.Run("TestWalkFlat", func(t *testing.T) {
		assertWalk(t, fs, path.Join(dir, "foo"), &WalkOptions{SkipSubDirs: true}, []string{"bar"}, []string{}, []string{})
	})

	t.Run("TestWalkError", func(t *testing.T) {
		errTest := errors.New("TestError")

		errors.Assert(t, errTest, fs.Walk(path.Join(dir, "foo"), func(dir string, f FileInfo, isRoot bool) errors.Error {
			if f.Name() == "test.txt" {
				// wait for inner directory to also test recursive error passing
				return errTest.Make()
			}
			return nil
		}, nil, nil, nil))

		errors.Assert(t, errTest, fs.Walk(path.Join(dir, "foo"), nil, func(dir string, f FileInfo, isRoot bool, skipDir *bool) errors.Error {
			return errTest.Make()
		}, nil, nil))

		errors.Assert(t, errTest, fs.Walk(path.Join(dir, "foo"), nil, nil, func(dir string, f FileInfo, isRoot bool) errors.Error {
			return errTest.Make()
		}, nil))
	})

	t.Run("TestWalkSkipDir", func(t *testing.T) {
		visitCount := 0
		visitExpected := []string{"bar"}
		errors.AssertNil(t, fs.Walk(path.Join(dir, "foo"), func(dir string, f FileInfo, isRoot bool) errors.Error {
			assert.Equal(t, visitExpected[visitCount], f.Name())
			visitCount++
			return nil
		}, func(dir string, f FileInfo, isRoot bool, skipDir *bool) errors.Error {
			*skipDir = true
			return nil
		}, nil, nil))
		assert.Equal(t, len(visitExpected), visitCount)
	})

	t.Run("TestMoveFile", func(t *testing.T) {
		errors.AssertNil(t, fs.Move(path.Join(dir, "foo/bar/test.txt"), path.Join(dir, "foo/test.txt")))
		assertNotExists(t, fs, path.Join(dir, "foo/bar/test.txt"))
		assertFileContent(t, fs, path.Join(dir, "foo/test.txt"), "foo bar")
	})

	t.Run("TestMoveDir", func(t *testing.T) {
		errors.AssertNil(t, fs.Move(path.Join(dir, "foo"), path.Join(dir, "asdf")))
		assertNotExists(t, fs, path.Join(dir, "foo"))
		assertIsDir(t, fs, path.Join(dir, "asdf/bar"))
		assertFileContent(t, fs, path.Join(dir, "asdf/test.txt"), "foo bar")
	})

	t.Run("TestMoveAll", func(t *testing.T) {
		errors.AssertNil(t, fs.MoveAll(path.Join(dir, "asdf"), path.Join(dir)))
		assertNotExists(t, fs, path.Join(dir, "asdf/bar"))
		assertNotExists(t, fs, path.Join(dir, "asdf/test.txt"))
		assertIsDir(t, fs, path.Join(dir, "bar"))
		assertFileContent(t, fs, path.Join(dir, "test.txt"), "foo bar")
	})

	t.Run("TestCopyFile", func(t *testing.T) {
		errors.AssertNil(t, fs.Copy(path.Join(dir, "test.txt"), path.Join(dir, "bar/test.txt")))
		assertIsFile(t, fs, path.Join(dir, "test.txt"))
		assertFileContent(t, fs, path.Join(dir, "bar/test.txt"), "foo bar")
	})

	t.Run("TestCopyDir", func(t *testing.T) {
		errors.AssertNil(t, fs.Copy(path.Join(dir, "bar"), path.Join(dir, "asdf/bar")))
		assertIsDir(t, fs, path.Join(dir, "bar"))
		assertFileContent(t, fs, path.Join(dir, "asdf/bar/test.txt"), "foo bar")
	})

	t.Run("TestWalkComplex", func(t *testing.T) {
		errors.AssertNil(t, fs.Move(path.Join(dir, "bar"), path.Join(dir, "asdf/test")))
		errors.AssertNil(t, fs.Move(path.Join(dir, "test.txt"), path.Join(dir, "asdf/file.txt")))

		dirCount := 0
		fileCount := 0
		size := int64(0)
		errors.AssertNil(t, fs.Walk(path.Join(dir, "asdf"), func(dir string, f FileInfo, isRoot bool) errors.Error {
			if f.IsDir() {
				dirCount++
			} else {
				fileCount++
				size += f.Size()
			}
			return nil
		}, nil, nil, nil))
		assert.Equal(t, 2, dirCount)
		assert.Equal(t, 3, fileCount)
		assert.Equal(t, int64(21), size)
	})

	t.Run("TestWalkFilesFirst", func(t *testing.T) {
		errors.AssertNil(t, fs.CreateDirectory(path.Join(dir, "foo2")))
		errors.AssertNil(t, fs.CreateDirectory(path.Join(dir, "foo2/bar")))
		errors.AssertNil(t, fs.WriteString(path.Join(dir, "foo2/test.txt"), "foo bar"))
		assertWalk(t, fs, path.Join(dir, "foo2"), &WalkOptions{VisitOrder: OrderFilesFirst}, []string{"test.txt", "bar"}, []string{"bar"}, []string{"bar"})
	})

	t.Run("TestWalkDirectoriesFirst", func(t *testing.T) {
		assertWalk(t, fs, path.Join(dir, "foo2"), &WalkOptions{VisitOrder: OrderDirectoriesFirst}, []string{"bar", "test.txt"}, []string{"bar"}, []string{"bar"})
	})

	t.Run("TestWalkLexicographicAsc", func(t *testing.T) {
		errors.AssertNil(t, fs.CreateDirectory(path.Join(dir, "foo2/cool")))
		errors.AssertNil(t, fs.CreateDirectory(path.Join(dir, "foo2/cool/sub")))
		errors.AssertNil(t, fs.CreateDirectory(path.Join(dir, "foo2/stuff")))
		errors.AssertNil(t, fs.WriteString(path.Join(dir, "foo2/better stuff.txt"), "foo bar"))
		assertWalk(t, fs, path.Join(dir, "foo2"), &WalkOptions{VisitOrder: OrderLexicographicAsc}, []string{"bar", "better stuff.txt", "cool", "sub", "stuff", "test.txt"}, []string{"bar", "cool", "sub", "stuff"}, []string{"bar", "sub", "cool", "stuff"})
	})

	t.Run("TestWalkLexicographicDesc", func(t *testing.T) {
		assertWalk(t, fs, path.Join(dir, "foo2"), &WalkOptions{VisitOrder: OrderLexicographicDesc}, []string{"test.txt", "stuff", "cool", "sub", "better stuff.txt", "bar"}, []string{"stuff", "cool", "sub", "bar"}, []string{"stuff", "sub", "cool", "bar"})
	})

	t.Run("TestWalkCompound", func(t *testing.T) {
		assertWalk(t, fs, path.Join(dir, "foo2"), &WalkOptions{VisitOrder: NewCompoundComparer(OrderFilesFirst, OrderLexicographicAsc)}, []string{"better stuff.txt", "test.txt", "bar", "cool", "sub", "stuff"}, []string{"bar", "cool", "sub", "stuff"}, []string{"bar", "sub", "cool", "stuff"})
	})
}

/* ############################################### */
/* ###               Test Helper               ### */
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

func assertWalk(t *testing.T, fs *FileSystem, path string, options *WalkOptions, visitExpected, enterExpected, leaveExpected []string) bool {
	visitRoot := false
	if options != nil && options.VisitRootDir {
		visitRoot = true
	}

	rootCallback := false
	if options != nil && options.EnterLeaveCallbacksForRoot {
		rootCallback = true
	}

	errAssertionFailed := errors.New("AssertionFailed")
	visitCount := 0
	enterCount := 0
	leaveCount := 0
	err := fs.Walk(path, func(dir string, f FileInfo, isRoot bool) errors.Error {
		if visitRoot && visitCount == 0 {
			if !assert.True(t, isRoot) {
				return errAssertionFailed.Make()
			}
		} else {
			if !assert.False(t, isRoot) {
				return errAssertionFailed.Make()
			}
		}
		if visitCount < len(visitExpected) && !assert.Equal(t, visitExpected[visitCount], f.Name()) {
			return errAssertionFailed.Make()
		}
		visitCount++
		return nil
	}, func(dir string, f FileInfo, isRoot bool, skipDir *bool) errors.Error {
		if rootCallback && enterCount == 0 {
			if !assert.True(t, isRoot) {
				return errAssertionFailed.Make()
			}
		} else {
			if !assert.False(t, isRoot) {
				return errAssertionFailed.Make()
			}
		}
		if enterCount < len(enterExpected) && !assert.Equal(t, enterExpected[enterCount], f.Name()) {
			return errAssertionFailed.Make()
		}
		enterCount++
		return nil
	}, func(dir string, f FileInfo, isRoot bool) errors.Error {
		if rootCallback && leaveCount == (len(leaveExpected)-1) {
			if !assert.True(t, isRoot) {
				return errAssertionFailed.Make()
			}
		} else {
			if !assert.False(t, isRoot) {
				return errAssertionFailed.Make()
			}
		}
		if leaveCount < len(leaveExpected) && !assert.Equal(t, leaveExpected[leaveCount], f.Name()) {
			return errAssertionFailed.Make()
		}
		leaveCount++
		return nil
	}, options)

	if errors.InstanceOf(err, errAssertionFailed) {
		return false
	}

	if !assert.Equal(t, len(visitExpected), visitCount) {
		return false
	}
	if !assert.Equal(t, len(enterExpected), enterCount) {
		return false
	}
	if !assert.Equal(t, len(leaveExpected), leaveCount) {
		return false
	}

	return true
}
