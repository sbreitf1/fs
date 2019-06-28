package fs

import (
	"io"

	"github.com/sbreitf1/errors"
)

var (
	ErrNotSupported       = errors.New("Operation %s is not supported by the file system")
	ErrFileNotExists      = errors.New("The file %q does not exist")
	ErrDirectoryNotExists = errors.New("The directory %q does not exist")
)

type ReadFileSystemDriver interface {
	Open(path string) (File, errors.Error)
}

type ReadWriteFileSystemDriver interface {
	ReadFileSystemDriver
}

type TempFileSystemDriver interface {
	ReadWriteFileSystemDriver
	GetTempFile(prefix string) (string, errors.Error)
}

type FileSystemDriver interface {
	TempFileSystemDriver
}

type FileInfo interface {
	Name() string
	Size() int64
}

type File interface {
	io.Reader
	io.Writer
	io.Closer
}

type FileSystem struct {
	driver                     interface{}
	rDriver                    ReadFileSystemDriver
	rwDriver                   ReadWriteFileSystemDriver
	tmpDriver                  TempFileSystemDriver
	canRead, canWrite, canTemp bool
}

func New(driver interface{}) *FileSystem {
	rDriver, rDriverOk := driver.(ReadFileSystemDriver)
	rwDriver, rwDriverOk := driver.(ReadWriteFileSystemDriver)
	tmpDriver, tmpDriverOk := driver.(TempFileSystemDriver)
	if !rDriverOk && !rwDriverOk && !tmpDriverOk {
		panic("fs.New expects valid File System Driver")
	}
	return &FileSystem{driver, rDriver, rwDriver, tmpDriver, rDriverOk, rwDriverOk, tmpDriverOk}
}

func NewLocal() *FileSystem {
	return New(NewLocalFileSystemDriver())
}

func (fs *FileSystem) CanRead() bool {
	return fs.canRead
}

func (fs *FileSystem) CanWrite() bool {
	return fs.canWrite
}

func (fs *FileSystem) CanTemp() bool {
	return fs.canTemp
}

func (fs *FileSystem) Open(path string) (File, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("Open").Make()
	}

	return fs.driver.(ReadFileSystemDriver).Open(path)
}
