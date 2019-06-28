package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var fs *Util
	assert.NotPanics(t, func() { fs = New() })
	assert.True(t, fs.CanRead(), "CanRead() returns false")
	assert.True(t, fs.CanWrite(), "CanRead() returns false")
	assert.True(t, fs.CanReadWrite(), "CanReadWrite() returns false")
	assert.True(t, fs.CanTemp(), "CanTemp() returns false")
	assert.True(t, fs.CanAll(), "CanAll() returns false")
}

func TestNewUtilInvalid(t *testing.T) {
	assert.Panics(t, func() { NewUtil("not a file system driver") })
}
