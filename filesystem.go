package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/sbreitf1/errors"
)

const (
	DefaultPathDelimiter = "/"
	DefaultLineDelimiter = "\n"
)

var (
	ErrNotSupported       = errors.New("Operation %s is not supported by the file system")
	ErrInvalidPath        = errors.New("Malformed path")
	ErrFileNotExists      = errors.New("The file %q does not exist")
	ErrDirectoryNotExists = errors.New("The directory %q does not exist")
)

type ReadFileSystemDriver interface {
	Exists(path string) (bool, errors.Error)
	IsFile(path string) (bool, errors.Error)
	IsDir(path string) (bool, errors.Error)

	Open(path string) (File, errors.Error)
}

type ReadWriteFileSystemDriver interface {
	ReadFileSystemDriver

	Create(path string) (File, errors.Error)
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

type Util struct {
	driver                     interface{}
	rDriver                    ReadFileSystemDriver
	rwDriver                   ReadWriteFileSystemDriver
	tmpDriver                  TempFileSystemDriver
	canRead, canWrite, canTemp bool
	lineSeparator              string
}

// New returns a new file system Util with local file system driver.
func New() *Util {
	driver := &LocalDriver{}
	return NewUtil(driver)
}

// NewUtil returns a new file system Util with the given file system driver.
func NewUtil(driver interface{}) *Util {
	rDriver, rDriverOk := driver.(ReadFileSystemDriver)
	rwDriver, rwDriverOk := driver.(ReadWriteFileSystemDriver)
	tmpDriver, tmpDriverOk := driver.(TempFileSystemDriver)
	if !rDriverOk && !rwDriverOk && !tmpDriverOk {
		panic(fmt.Sprintf("fs.New expects valid File System Driver but got %T instead", driver))
	}
	return &Util{driver, rDriver, rwDriver, tmpDriver, rDriverOk, rwDriverOk, tmpDriverOk, DefaultLineDelimiter}
}

func (fs *Util) CanRead() bool {
	return fs.canRead
}

func (fs *Util) CanWrite() bool {
	return fs.canWrite
}

func (fs *Util) CanReadWrite() bool {
	return fs.canRead && fs.canWrite
}

func (fs *Util) CanTemp() bool {
	return fs.canTemp
}

func (fs *Util) CanAll() bool {
	return fs.canRead && fs.canWrite && fs.canTemp
}

// Exists returns true, if the given path is a file or directory.
func (fs *Util) Exists(path string) (bool, errors.Error) {
	if !fs.canRead {
		return false, ErrNotSupported.Args("Exists").Make()
	}

	return fs.rDriver.Exists(path)
}

// IsFile returns true, if the given path is a file.
func (fs *Util) IsFile(path string) (bool, errors.Error) {
	if !fs.canRead {
		return false, ErrNotSupported.Args("IsFile").Make()
	}

	return fs.rDriver.IsFile(path)
}

// IsDir returns true, if the given path is a directory.
func (fs *Util) IsDir(path string) (bool, errors.Error) {
	if !fs.canRead {
		return false, ErrNotSupported.Args("IsDir").Make()
	}

	return fs.rDriver.IsDir(path)
}

func (fs *Util) Open(path string) (File, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("Open").Make()
	}

	return fs.rDriver.Open(path)
}

func (fs *Util) ReadBytes(path string) ([]byte, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("ReadBytes").Make()
	}

	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, readErr := ioutil.ReadAll(f)
	if readErr != nil {
		return nil, errors.Wrap(readErr).Expand("Failed to read file")
	}

	return data, nil
}

func (fs *Util) ReadString(path string) (string, errors.Error) {
	if !fs.canRead {
		return "", ErrNotSupported.Args("ReadString").Make()
	}

	data, err := fs.ReadBytes(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (fs *Util) ReadLines(path string) ([]string, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("ReadLines").Make()
	}

	str, err := fs.ReadString(path)
	if err != nil {
		return nil, err
	}

	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "\n", -1)
	return strings.Split(str, "\n"), nil
}

func (fs *Util) Create(path string) (File, errors.Error) {
	if !fs.canWrite {
		return nil, ErrNotSupported.Args("Create").Make()
	}

	return fs.rwDriver.Create(path)
}

func (fs *Util) WriteBytes(path string, content []byte) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("WriteBytes").Make()
	}

	f, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	//TODO check all bytes written
	if _, err := f.Write(content); err != nil {
		return errors.Wrap(err).Expand("Failed to write file")
	}
	return nil
}

func (fs *Util) WriteString(path, content string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("WriteString").Make()
	}

	return fs.WriteBytes(path, []byte(content))
}

func (fs *Util) WriteLines(path string, lines []string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("WriteLines").Make()
	}

	return fs.WriteBytes(path, []byte(strings.Join(lines, fs.lineSeparator)))
}
