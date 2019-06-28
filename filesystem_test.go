package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInvalid(t *testing.T) {
	assert.Panics(t, func() { New("not a file system driver") })
}

func TestNewLocal(t *testing.T) {
	var fs *FileSystem
	assert.NotPanics(t, func() { fs = NewLocal() })
	assert.True(t, fs.CanRead(), "CanRead() returns false")
	assert.True(t, fs.CanWrite(), "CanRead() returns false")
	assert.True(t, fs.CanTemp(), "CanTemp() returns false")
}
