package path

import (
	"os"
	"testing"

	"github.com/sbreitf1/errors"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
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
