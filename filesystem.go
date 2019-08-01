package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sbreitf1/errors"
)

const (
	// DefaultLineDelimiter denotes the character or characther sequence to separate lines in text files.
	DefaultLineDelimiter = "\n"
)

var (
	// ErrFileSystem is a generic file system related error.
	ErrFileSystem = errors.New("A file system error occured")
	// ErrNotSupported is returned when using a function that is not supported.
	ErrNotSupported = errors.New("Operation %s is not supported by the file system")
	// ErrInvalidPath occurs when using malformed paths.
	ErrInvalidPath = errors.New("Malformed path")
	// ErrFileNotExists occurs when an action failed because of a missing file.
	ErrFileNotExists = errors.New("The file %q does not exist")
	// ErrDirectoryNotExists occurs when an action failed because of a missing directory.
	ErrDirectoryNotExists = errors.New("The directory %q does not exist")
	// ErrAccessDenied denotes an error caused by insufficient privileges.
	ErrAccessDenied = errors.New("Acces to %q denied")
)

// ReadFileSystemDriver describes functionality to read from a file system.
type ReadFileSystemDriver interface {
	Exists(path string) (bool, errors.Error)
	IsFile(path string) (bool, errors.Error)
	IsDir(path string) (bool, errors.Error)

	ReadDir(path string) ([]FileInfo, errors.Error)

	OpenFile(path string) (File, errors.Error)
}

// ReadWriteFileSystemDriver describes functionality to write to a file system.
type ReadWriteFileSystemDriver interface {
	ReadFileSystemDriver

	CreateFile(path string) (File, errors.Error)

	DeleteFile(path string) errors.Error
	DeleteDirectory(path string, recursive bool) errors.Error

	MoveFile(src, dst string) errors.Error
	MoveDir(src, dst string) errors.Error
}

// TempFileSystemDriver describes functionality to create temporary files and directories on a file system.
type TempFileSystemDriver interface {
	ReadWriteFileSystemDriver

	GetTempFile(pattern string) (string, errors.Error)
	GetTempDir(prefix string) (string, errors.Error)
}

// FileSystemDriver describes a complete file system function set.
type FileSystemDriver interface {
	TempFileSystemDriver
}

// FileInfo contains meta information for a file.
type FileInfo interface {
	Name() string
	Size() int64
	IsDir() bool
}

// File is the instance object for an opened file.
type File interface {
	io.Reader
	io.Writer
	io.Closer
}

// FileSystem offers advanced functionality based on a file system driver.
type FileSystem struct {
	driver                     interface{}
	rDriver                    ReadFileSystemDriver
	rwDriver                   ReadWriteFileSystemDriver
	tmpDriver                  TempFileSystemDriver
	canRead, canWrite, canTemp bool
	lineSeparator              string
}

// New returns a new file system with local file system driver.
func New() *FileSystem {
	return NewWithDriver(&LocalDriver{})
}

// NewWithDriver returns a new file system using the given file system driver.
func NewWithDriver(driver interface{}) *FileSystem {
	rDriver, rDriverOk := driver.(ReadFileSystemDriver)
	rwDriver, rwDriverOk := driver.(ReadWriteFileSystemDriver)
	tmpDriver, tmpDriverOk := driver.(TempFileSystemDriver)
	if !rDriverOk && !rwDriverOk && !tmpDriverOk {
		panic(fmt.Sprintf("fs.New expects valid File System Driver but got %T instead", driver))
	}
	return &FileSystem{driver, rDriver, rwDriver, tmpDriver, rDriverOk, rwDriverOk, tmpDriverOk, DefaultLineDelimiter}
}

// CanRead returns true when the file system can perform read operations.
func (fs *FileSystem) CanRead() bool {
	return fs.canRead
}

// CanWrite returns true when the file system can perform write operations.
func (fs *FileSystem) CanWrite() bool {
	return fs.canWrite
}

// CanReadWrite returns true when the file system can perform both read and write operations.
func (fs *FileSystem) CanReadWrite() bool {
	return fs.canRead && fs.canWrite
}

// CanTemp returns true when the file system can create temporary files and directories.
func (fs *FileSystem) CanTemp() bool {
	return fs.canTemp
}

// CanAll returns true when the file system offers complete functionality.
func (fs *FileSystem) CanAll() bool {
	return fs.canRead && fs.canWrite && fs.canTemp
}

/* ############################################### */
/* ###               Read Access               ### */
/* ############################################### */

// Exists returns true, if the given path is a file or directory.
func (fs *FileSystem) Exists(path string) (bool, errors.Error) {
	if !fs.canRead {
		return false, ErrNotSupported.Args("Exists").Make()
	}

	return fs.rDriver.Exists(path)
}

// IsFile returns true, if the given path is a file.
func (fs *FileSystem) IsFile(path string) (bool, errors.Error) {
	if !fs.canRead {
		return false, ErrNotSupported.Args("IsFile").Make()
	}

	return fs.rDriver.IsFile(path)
}

// IsDir returns true, if the given path is a directory.
func (fs *FileSystem) IsDir(path string) (bool, errors.Error) {
	if !fs.canRead {
		return false, ErrNotSupported.Args("IsDir").Make()
	}

	return fs.rDriver.IsDir(path)
}

// ReadDir returns all files and directories contained in a directory.
func (fs *FileSystem) ReadDir(path string) ([]FileInfo, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("ReadDir").Make()
	}

	return fs.rDriver.ReadDir(path)
}

// OpenFile opens a file instance and returns the handle.
func (fs *FileSystem) OpenFile(path string) (File, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("Open").Make()
	}

	return fs.rDriver.OpenFile(path)
}

// ReadBytes returns all bytes contained in a file.
func (fs *FileSystem) ReadBytes(path string) ([]byte, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("ReadBytes").Make()
	}

	f, err := fs.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, readErr := ioutil.ReadAll(f)
	if readErr != nil {
		return nil, ErrFileSystem.Msg("Failed to read file").Make().Cause(readErr)
	}

	return data, nil
}

// ReadString returns the file content as string.
func (fs *FileSystem) ReadString(path string) (string, errors.Error) {
	if !fs.canRead {
		return "", ErrNotSupported.Args("ReadString").Make()
	}

	data, err := fs.ReadBytes(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadLines returns all files contained in a text file.
func (fs *FileSystem) ReadLines(path string) ([]string, errors.Error) {
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

/* ############################################### */
/* ###              Write Access               ### */
/* ############################################### */

// CreateFile a new file (or truncate an existing) and return the file instance handle.
func (fs *FileSystem) CreateFile(path string) (File, errors.Error) {
	if !fs.canWrite {
		return nil, ErrNotSupported.Args("Create").Make()
	}

	return fs.rwDriver.CreateFile(path)
}

// WriteBytes writes all bytes to a file.
func (fs *FileSystem) WriteBytes(path string, content []byte) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("WriteBytes").Make()
	}

	f, err := fs.CreateFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return ErrFileSystem.Msg("Failed to write file").Make().Cause(err)
	}
	return nil
}

// WriteString writes a string to a file.
func (fs *FileSystem) WriteString(path, content string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("WriteString").Make()
	}

	return fs.WriteBytes(path, []byte(content))
}

// WriteLines writes all lines to a file using the default line delimiter.
func (fs *FileSystem) WriteLines(path string, lines []string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("WriteLines").Make()
	}

	return fs.WriteBytes(path, []byte(strings.Join(lines, fs.lineSeparator)))
}

// CopyFile clones a file and overwrites the existing one.
func (fs *FileSystem) CopyFile(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("CopyFile").Make()
	}

	reader, err := fs.OpenFile(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := fs.CreateFile(dst)
	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		return ErrFileSystem.Msg("Failed to copy file").Make().Cause(err)
	}
	return nil
}

// CleanDir removes all files and directories from a directory.
func (fs *FileSystem) CleanDir(path string) errors.Error {
	return nil
}

/* ############################################### */
/* ###               Temp Files                ### */
/* ############################################### */

/* ############################################### */
/* ###                Contexts                 ### */
/* ############################################### */

// WithTempFile creates a temporary file and deletes it when f returns.
func (fs *FileSystem) WithTempFile(pattern string, f func(tmpFile string) errors.Error) errors.Error {
	tmpFile, err := fs.tmpDriver.GetTempFile(pattern)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	return f(tmpFile)
}

// WithTempDir creates a temporary directory and deletes it when f returns.
func (fs *FileSystem) WithTempDir(prefix string, f func(tmpDir string) errors.Error) errors.Error {
	tmpDir, err := fs.tmpDriver.GetTempDir(prefix)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	return f(tmpDir)
}
