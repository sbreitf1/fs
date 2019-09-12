package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sbreitf1/fs/path"

	"github.com/sbreitf1/errors"
)

const (
	// DefaultLineDelimiter denotes the character or characther sequence to separate lines in text files.
	DefaultLineDelimiter = "\n"
)

var (
	// Err is a generic file system related error.
	Err = errors.New("A file system error occured")
	// ErrNotSupported is returned when using a function that is not supported.
	ErrNotSupported = errors.New("Operation %s is not supported by the file system")
	// ErrNotExists occurs when an action failed becaus of a missing element of unknown type(file or directory).
	ErrNotExists = errors.New("The path %q does not exist")
	// ErrFileNotExists occurs when an action failed because of a missing file.
	ErrFileNotExists = errors.New("The file %q does not exist")
	// ErrDirectoryNotExists occurs when an action failed because of a missing directory.
	ErrDirectoryNotExists = errors.New("The directory %q does not exist")
	// ErrAccessDenied denotes an error caused by insufficient privileges.
	ErrAccessDenied = errors.New("Access to %q denied")
	// ErrNotEmpty occurs when trying to delete a non-empty directory without recursive flag.
	ErrNotEmpty = errors.New("The directory is not empty")
)

// NavigationFileSystemDriver describes functionality to list files and directories but does not allow access to file content.
type NavigationFileSystemDriver interface {
	Exists(path string) (bool, errors.Error)
	IsFile(path string) (bool, errors.Error)
	IsDir(path string) (bool, errors.Error)

	Stat(path string) (FileInfo, errors.Error)
	ReadDir(path string) ([]FileInfo, errors.Error)
}

// ReadFileSystemDriver describes functionality to read from a file system.
type ReadFileSystemDriver interface {
	NavigationFileSystemDriver

	OpenFile(path string, flags OpenFlags) (File, errors.Error)
}

// ReadWriteFileSystemDriver describes functionality to write to a file system.
type ReadWriteFileSystemDriver interface {
	ReadFileSystemDriver

	CreateDirectory(path string) errors.Error

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

// OpenFlags specifies information on how to open a file.
type OpenFlags int

const (
	// OpenReadOnly denotes opening a file using read-only access.
	OpenReadOnly OpenFlags = OpenFlags(os.O_RDONLY)
	// OpenWriteOnly denotes opening a file using write-only access.
	OpenWriteOnly OpenFlags = OpenFlags(os.O_WRONLY)
	// OpenReadWrite denotes opening a file using read-write access.
	OpenReadWrite OpenFlags = OpenFlags(os.O_RDWR)
)

// Access returns only the access flag.
func (flag OpenFlags) Access() OpenFlags {
	mask := int(OpenReadOnly) | int(OpenWriteOnly) | int(OpenReadWrite)
	return OpenFlags(int(flag) & mask)
}

// IsRead returns whether the given flags require read access.
func (flag OpenFlags) IsRead() bool {
	access := flag.Access()
	return (access == OpenReadOnly || access == OpenReadWrite)
}

// IsWrite returns whether the given flags require write access.
func (flag OpenFlags) IsWrite() bool {
	access := flag.Access()
	return (access == OpenWriteOnly || access == OpenReadWrite)
}

// Append opens the file for appending.
func (flag OpenFlags) Append() OpenFlags {
	return OpenFlags(int(flag) | os.O_APPEND)
}

// Create creates the file if it does not exist.
func (flag OpenFlags) Create() OpenFlags {
	return OpenFlags(int(flag) | os.O_CREATE)
}

// Exclusive opens the file for appending.
func (flag OpenFlags) Exclusive() OpenFlags {
	return OpenFlags(int(flag) | os.O_EXCL)
}

// Sync opens the file for appending.
func (flag OpenFlags) Sync() OpenFlags {
	return OpenFlags(int(flag) | os.O_SYNC)
}

// Truncate opens the file for appending.
func (flag OpenFlags) Truncate() OpenFlags {
	return OpenFlags(int(flag) | os.O_TRUNC)
}

// FileSystem offers advanced functionality based on a file system driver.
type FileSystem struct {
	navDriver                               NavigationFileSystemDriver
	rDriver                                 ReadFileSystemDriver
	rwDriver                                ReadWriteFileSystemDriver
	tmpDriver                               TempFileSystemDriver
	canNavigate, canRead, canWrite, canTemp bool
	LineSeparator                           string
}

// New returns a new file system with local file system driver.
func New() *FileSystem {
	return NewWithDriver(&LocalDriver{})
}

// NewWithDriver returns a new file system using the given file system driver. The given driver must implement at least one of the file system driver interfaces.
func NewWithDriver(driver interface{}) *FileSystem {
	navDriver, navDriverOk := driver.(NavigationFileSystemDriver)
	rDriver, rDriverOk := driver.(ReadFileSystemDriver)
	rwDriver, rwDriverOk := driver.(ReadWriteFileSystemDriver)
	tmpDriver, tmpDriverOk := driver.(TempFileSystemDriver)
	if !rDriverOk && !rwDriverOk && !tmpDriverOk {
		//TODO show message if driver is not passed as pointer
		panic(fmt.Sprintf("fs.New expects valid File System Driver but got %T instead", driver))
	}
	return &FileSystem{navDriver, rDriver, rwDriver, tmpDriver, navDriverOk, rDriverOk, rwDriverOk, tmpDriverOk, DefaultLineDelimiter}
}

// CanNavigate returns true when the file system allows to list files and directories.
func (fs *FileSystem) CanNavigate() bool {
	return fs.canNavigate
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
	return fs.canNavigate && fs.canRead && fs.canWrite && fs.canTemp
}

/* ############################################### */
/* ###            Navigation Access            ### */
/* ############################################### */

// Exists returns true, if the given path is a file or directory.
func (fs *FileSystem) Exists(path string) (bool, errors.Error) {
	if !fs.canNavigate {
		return false, ErrNotSupported.Args("Exists").Make()
	}

	return fs.navDriver.Exists(path)
}

// IsFile returns true, if the given path is a file.
func (fs *FileSystem) IsFile(path string) (bool, errors.Error) {
	if !fs.canNavigate {
		return false, ErrNotSupported.Args("IsFile").Make()
	}

	return fs.navDriver.IsFile(path)
}

// IsDir returns true, if the given path is a directory.
func (fs *FileSystem) IsDir(path string) (bool, errors.Error) {
	if !fs.canNavigate {
		return false, ErrNotSupported.Args("IsDir").Make()
	}

	return fs.navDriver.IsDir(path)
}

// Stat returns file or directory stats for a given path.
func (fs *FileSystem) Stat(path string) (FileInfo, errors.Error) {
	if !fs.canNavigate {
		return nil, ErrNotSupported.Args("Stat").Make()
	}

	return fs.navDriver.Stat(path)
}

// ReadDir returns all files and directories contained in a directory.
func (fs *FileSystem) ReadDir(path string) ([]FileInfo, errors.Error) {
	if !fs.canNavigate {
		return nil, ErrNotSupported.Args("ReadDir").Make()
	}

	return fs.navDriver.ReadDir(path)
}

// EnterDirHandler is called by Walk before a directory is entered. If skipDir is set to true, the directory will not be visited.
type EnterDirHandler func(dir string, f FileInfo, skipDir *bool) errors.Error

// VisitFileHandler is called by Walk for every file that is found recursively.
type VisitFileHandler func(dir string, f FileInfo) errors.Error

// LeaveDirHandler is called by Walk after all elements inside a directory have been processed.
type LeaveDirHandler func(dir string, f FileInfo) errors.Error

// WalkOptions can be used to specify the behavior of Walk like visit order and search strategy.
type WalkOptions struct {
	// SkipSubDirs denotes whether the directory is traversed recursively or not.
	SkipSubDirs bool

	//TODO walk options
}

// Walk calls the corresponding callback functions for ever file and directory contained in dir recursively.
func (fs *FileSystem) Walk(dir string, visitFileHandler VisitFileHandler, enterDirHandler EnterDirHandler, leaveDirHandler LeaveDirHandler, options *WalkOptions) errors.Error {
	if !fs.canNavigate {
		return ErrNotSupported.Args("Walk").Make()
	}

	if options == nil {
		options = &WalkOptions{}
	}

	return fs.walk(dir, visitFileHandler, enterDirHandler, leaveDirHandler, options)
}

func (fs *FileSystem) walk(dir string, visitFileHandler VisitFileHandler, enterDirHandler EnterDirHandler, leaveDirHandler LeaveDirHandler, options *WalkOptions) errors.Error {
	files, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if visitFileHandler != nil {
			if err := visitFileHandler(dir, f); err != nil {
				return err
			}
		}

		if !options.SkipSubDirs && f.IsDir() {
			if enterDirHandler != nil {
				skipDir := false
				if err := enterDirHandler(dir, f, &skipDir); err != nil {
					return err
				}
				if skipDir {
					return nil
				}
			}

			if err := fs.walk(path.Join(dir, f.Name()), visitFileHandler, enterDirHandler, leaveDirHandler, options); err != nil {
				return err
			}

			if leaveDirHandler != nil {
				if err := leaveDirHandler(dir, f); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/* ############################################### */
/* ###               Read Access               ### */
/* ############################################### */

// Open opens a file instance for reading and returns the handle.
func (fs *FileSystem) Open(path string) (File, errors.Error) {
	if !fs.canRead {
		return nil, ErrNotSupported.Args("Open").Make()
	}

	return fs.rDriver.OpenFile(path, OpenReadOnly)
}

// OpenFile opens a general purpose file instance based on flags and returns the handle.
func (fs *FileSystem) OpenFile(path string, flags OpenFlags) (File, errors.Error) {
	if flags.IsRead() && !fs.canRead {
		return nil, ErrNotSupported.Args("OpenFile (read)").Make()
	}
	if flags.IsWrite() && !fs.canWrite {
		return nil, ErrNotSupported.Args("OpenFile (write)").Make()
	}

	return fs.rDriver.OpenFile(path, flags)
}

// ReadBytes returns all bytes contained in a file.
func (fs *FileSystem) ReadBytes(path string) ([]byte, errors.Error) {
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
		return nil, Err.Msg("Failed to read file").Make().Cause(readErr)
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
		return nil, ErrNotSupported.Args("CreateFile").Make()
	}

	return fs.rwDriver.OpenFile(path, OpenReadWrite.Create().Truncate())
}

// CreateDirectory creates a new directory and all parent directories if they do not exist.
func (fs *FileSystem) CreateDirectory(path string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("CreateDirectory").Make()
	}

	return fs.rwDriver.CreateDirectory(path)
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
		return Err.Msg("Failed to write file").Make().Cause(err)
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

	return fs.WriteBytes(path, []byte(strings.Join(lines, fs.LineSeparator)))
}

// DeleteFile deletes a file.
func (fs *FileSystem) DeleteFile(path string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("DeleteFile").Make()
	}

	return fs.rwDriver.DeleteFile(path)
}

// DeleteDirectory deletes an empty directory. If recursive is set, all contained items will be deleted first.
func (fs *FileSystem) DeleteDirectory(path string, recursive bool) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("DeleteDirectory").Make()
	}

	return fs.rwDriver.DeleteDirectory(path, recursive)
}

// CleanDir removes all files and directories from a directory.
func (fs *FileSystem) CleanDir(dir string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("CleanDir").Make()
	}

	files, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if err := fs.DeleteDirectory(path.Join(dir, f.Name()), true); err != nil {
				return err
			}
		} else {
			if err := fs.DeleteFile(path.Join(dir, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

/* ############################################### */
/* ###             Move and Copy               ### */
/* ############################################### */

// Move moves a file or directory to a new location. If the target already exists, it must be the same element type (file or directory) to be overwritten.
func (fs *FileSystem) Move(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("Move").Make()
	}

	isFile, err := fs.IsFile(src)
	if err != nil {
		return err
	}
	if isFile {
		return fs.MoveFile(src, dst)
	}

	isDir, err := fs.IsDir(src)
	if err != nil {
		return err
	}
	if isDir {
		return fs.MoveDir(src, dst)
	}

	return ErrNotExists.Args(src).Make()
}

// MoveFile moves a file to a new location.
func (fs *FileSystem) MoveFile(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("MoveFile").Make()
	}

	return fs.rwDriver.MoveFile(src, dst)
}

// MoveDir moves a directory to a new location.
func (fs *FileSystem) MoveDir(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("MoveDir").Make()
	}

	return fs.rwDriver.MoveDir(src, dst)
}

// MoveAll moves all files and directories contained in src to dst.
func (fs *FileSystem) MoveAll(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("MoveAll").Make()
	}

	files, err := fs.rDriver.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if err := fs.MoveDir(path.Join(src, f.Name()), path.Join(dst, f.Name())); err != nil {
				return err
			}
		} else {
			if err := fs.MoveFile(path.Join(src, f.Name()), path.Join(dst, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

//TODO MoveDir with callback before overwrite (cancel/skip/overwrite/rename) -> maybe replace existing MoveDir method?
// -> specify default handlers for cancel / skip / overwrite and rename by adding a number

// Copy clone a file or directory to the target. If the target already exists, it must be the same element type (file or directory) to be overwritten.
func (fs *FileSystem) Copy(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("Copy").Make()
	}

	isFile, err := fs.IsFile(src)
	if err != nil {
		return err
	}
	if isFile {
		return fs.CopyFile(src, dst)
	}

	isDir, err := fs.IsDir(src)
	if err != nil {
		return err
	}
	if isDir {
		return fs.CopyDir(src, dst)
	}

	return ErrNotExists.Args(src).Make()
}

// CopyFile clones a file and overwrites the existing one.
func (fs *FileSystem) CopyFile(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("CopyFile").Make()
	}

	reader, err := fs.Open(src)
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
		return Err.Msg("Failed to copy file").Make().Cause(err)
	}
	return nil
}

// CopyDir recursively clones a directory overwriting all existing files.
func (fs *FileSystem) CopyDir(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("CopyDir").Make()
	}

	if err := fs.rwDriver.CreateDirectory(dst); err != nil {
		return err
	}

	files, err := fs.rDriver.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if err := fs.CopyDir(path.Join(src, f.Name()), path.Join(dst, f.Name())); err != nil {
				return err
			}
		} else {
			if err := fs.CopyFile(path.Join(src, f.Name()), path.Join(dst, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyAll copies all files and directories contained in src to dst.
func (fs *FileSystem) CopyAll(src, dst string) errors.Error {
	if !fs.canWrite {
		return ErrNotSupported.Args("CopyAll").Make()
	}

	files, err := fs.rDriver.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if err := fs.CopyDir(path.Join(src, f.Name()), path.Join(dst, f.Name())); err != nil {
				return err
			}
		} else {
			if err := fs.CopyFile(path.Join(src, f.Name()), path.Join(dst, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

//TODO CopyDir with callback before overwrite (cancel/skip/overwrite/rename)

/* ############################################### */
/* ###               Temp Files                ### */
/* ############################################### */

// GetTempFile returns the path to an empty temporary file.
func (fs *FileSystem) GetTempFile(pattern string) (string, errors.Error) {
	if !fs.canTemp {
		return "", ErrNotSupported.Args("GetTempFile").Make()
	}

	return fs.tmpDriver.GetTempFile(pattern)
}

// GetTempDir returns the path to an empty temporary dir.
func (fs *FileSystem) GetTempDir(prefix string) (string, errors.Error) {
	if !fs.canTemp {
		return "", ErrNotSupported.Args("GetTempDir").Make()
	}

	return fs.tmpDriver.GetTempDir(prefix)
}

/* ############################################### */
/* ###                Contexts                 ### */
/* ############################################### */

// WithTempFile creates a temporary file and deletes it when f returns.
func (fs *FileSystem) WithTempFile(pattern string, f func(tmpFile string) errors.Error) errors.Error {
	if !fs.canTemp {
		return ErrNotSupported.Args("WithTempFile").Make()
	}

	tmpFile, err := fs.tmpDriver.GetTempFile(pattern)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	return f(tmpFile)
}

// WithTempDir creates a temporary directory and deletes it when f returns.
func (fs *FileSystem) WithTempDir(prefix string, f func(tmpDir string) errors.Error) errors.Error {
	if !fs.canTemp {
		return ErrNotSupported.Args("WithTempDir").Make()
	}

	tmpDir, err := fs.tmpDriver.GetTempDir(prefix)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	return f(tmpDir)
}
