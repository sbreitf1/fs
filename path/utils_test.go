package path

import (
	"os"
	"testing"

	"github.com/sbreitf1/errors"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	assert.Equal(t, "", Join(""))
	assert.Equal(t, "/", Join("/"))
	assert.Equal(t, "/", Join("", "/"))
	assert.Equal(t, "/foobar", Join("", "/foobar"))
	assert.Equal(t, "/bar", Join("", "/bar"))
	assert.Equal(t, "/foo/bar", Join("/foo", "", "bar"))
	assert.Equal(t, "/foo/bar", Join("/foo", "/", "bar"))
	assert.Equal(t, "/foo/bar", Join("/foo/", "/", "/bar"))
	assert.Equal(t, "/foo/bar", Join("/foo", "bar"))
	assert.Equal(t, "/foo/bar", Join("/foo/", "bar"))
	assert.Equal(t, "/foo/bar", Join("/foo", "/bar"))
	assert.Equal(t, "/foo/bar", Join("/foo/", "/bar"))
	assert.Equal(t, "foo/bar/test", Join("foo/bar", "test"))
	assert.Equal(t, "foo/bar/test", Join("foo", "bar", "test"))
}

func TestBase(t *testing.T) {
	assert.Equal(t, ".", Base(""))
	assert.Equal(t, "bar", Base("/foo/bar"))
	assert.Equal(t, "foo", Base("foo"))
	assert.Equal(t, "foo", Base("/foo"))
	assert.Equal(t, "foo", Base("/foo/"))
	assert.Equal(t, "bar", Base("foo/bar"))
	assert.Equal(t, "test.jpg", Base("/foo/bar/test.jpg"))
}

func TestBaseNoExt(t *testing.T) {
	assert.Equal(t, ".", BaseNoExt(""))
	assert.Equal(t, "bar", BaseNoExt("/foo/bar"))
	assert.Equal(t, "test", BaseNoExt("/foo/bar/test.jpg"))
	assert.Equal(t, "archive.tar", BaseNoExt("/foo/bar/archive.tar.gz"))
}

func TestDir(t *testing.T) {
	assert.Equal(t, "/foo", Dir("/foo/bar"))
	assert.Equal(t, ".", Dir("foo"))
	assert.Equal(t, "/", Dir("/foo"))
	assert.Equal(t, "/foo", Dir("/foo/"))
	assert.Equal(t, "foo", Dir("foo/bar"))
	assert.Equal(t, "/foo/bar", Dir("/foo/bar/test.jpg"))
}

func TestClean(t *testing.T) {
	assert.Equal(t, "foo", Clean("foo"))
	assert.Equal(t, "/foo", Clean("/foo"))
	assert.Equal(t, "foo", Clean("./foo"))
	assert.Equal(t, "../foo", Clean("../foo"))
	assert.Equal(t, "bar", Clean("foo/../bar"))
	assert.Equal(t, "foo/bar", Clean("foo/./bar"))
	assert.Equal(t, "../bar", Clean("foo/../../bar"))
	assert.Equal(t, "/bar", Clean("/foo/../../bar"))
	assert.Equal(t, "/", Clean("/foo/../.."))
}

func TestIsAbs(t *testing.T) {
	assert.True(t, IsAbs("/foo"))
	assert.False(t, IsAbs("foo"))
	assert.False(t, IsAbs("./foo"))
	assert.False(t, IsAbs("../foo"))
	assert.True(t, IsAbs("/foo/bar/.."))
	assert.True(t, IsAbs("/foo/./bar"))
	assert.True(t, IsAbs("/foo/.."))
	assert.True(t, IsAbs("/foo/../.."))
}

func TestAbs(t *testing.T) {
	var abs string
	var err error

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	abs, err = Abs("test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, Join(wd, "test.txt"), abs)

	abs, err = Abs("./test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, Join(wd, "test.txt"), abs)

	abs, err = Abs("../test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, Join(Dir(wd), "test.txt"), abs)

	abs, err = Abs("/foo/bar/../test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/foo/test.txt", abs)

	abs, err = Abs("/foo/bar/..")
	errors.AssertNil(t, err)
	assert.Equal(t, "/foo", abs)

	abs, err = Abs("/foo/bar/../")
	errors.AssertNil(t, err)
	assert.Equal(t, "/foo", abs)

	abs, err = Abs("/./foo")
	errors.AssertNil(t, err)
	assert.Equal(t, "/foo", abs)

	abs, err = Abs("/../../../foo")
	errors.AssertNil(t, err)
	assert.Equal(t, "/foo", abs)
}

func TestAbsIn(t *testing.T) {
	var abs string
	var err error

	abs, err = AbsIn("var/blub", "test.txt")
	errors.Assert(t, Err, err)

	abs, err = AbsIn("/var/blub", "")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub", abs)

	abs, err = AbsIn("/var/blub", "/")
	errors.AssertNil(t, err)
	assert.Equal(t, "/", abs)

	abs, err = AbsIn("/var/blub", "./")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub", abs)

	abs, err = AbsIn("/var/blub", "test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsIn("/var/blub", "/test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/test.txt", abs)

	abs, err = AbsIn("/var/blub", "./test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsIn("/var/blub", "foo/../test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsIn("/var/blub", "../test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/test.txt", abs)

	abs, err = AbsIn("/var/blub", "../blubber")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blubber", abs)

	abs, err = AbsIn("/var/blub", "../../../../../blubber")
	errors.AssertNil(t, err)
	assert.Equal(t, "/blubber", abs)
}

func TestAbsRoot(t *testing.T) {
	var abs string
	var err error

	abs, err = AbsRoot("var/blub", "test.txt")
	errors.Assert(t, Err, err)

	abs, err = AbsRoot("/var/blub", "")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub", abs)

	abs, err = AbsRoot("/var/blub", "/")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub", abs)

	abs, err = AbsRoot("/var/blub", "./")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub", abs)

	abs, err = AbsRoot("/var/blub", "test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsRoot("/var/blub", "/test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsRoot("/var/blub", "./test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsRoot("/var/blub", "foo/../test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsRoot("/var/blub", "../test.txt")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/test.txt", abs)

	abs, err = AbsRoot("/var/blub", "../blubber")
	errors.AssertNil(t, err)
	assert.Equal(t, "/var/blub/blubber", abs)
}

func TestIsIn(t *testing.T) {
	var in bool
	var err errors.Error

	in, err = IsIn("/", "/")
	errors.AssertNil(t, err)
	assert.True(t, in)

	in, err = IsIn("/var/blub", "/var")
	errors.AssertNil(t, err)
	assert.True(t, in)

	in, err = IsIn("/usr/bin", "/var")
	errors.AssertNil(t, err)
	assert.False(t, in)

	in, err = IsIn("/usr/bin", "/usr/binner")
	errors.AssertNil(t, err)
	assert.False(t, in)

	in, err = IsIn("/", "/usr/bin")
	errors.AssertNil(t, err)
	assert.False(t, in)

	in, err = IsIn("", "/")
	errors.Assert(t, Err, err)

	in, err = IsIn("/", "")
	errors.Assert(t, Err, err)
}

func TestExt(t *testing.T) {
	assert.Equal(t, ".jpg", Ext("image.jpg"))
	assert.Equal(t, ".png", Ext("/usr/image.png"))
	assert.Equal(t, "", Ext("/home/just-a-file"))
	assert.Equal(t, ".gz", Ext("/home/Downloads/archive.tar.gz"))
}

func TestNoExt(t *testing.T) {
	assert.Equal(t, "image", NoExt("image.jpg"))
	assert.Equal(t, "/usr/image", NoExt("/usr/image.png"))
	assert.Equal(t, "/home/just-a-file", NoExt("/home/just-a-file"))
	assert.Equal(t, "/home/Downloads/archive.tar", NoExt("/home/Downloads/archive.tar.gz"))
}
