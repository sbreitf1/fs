package fs

import (
	"github.com/sbreitf1/errors"
)

var (
	// DefaultFileSystem denots the file system used for all default accessors.
	DefaultFileSystem *FileSystem
)

func init() {
	DefaultFileSystem = New()
}

// Exists returns true, if the given path is a file or directory.
func Exists(path string) (bool, errors.Error) {
	return DefaultFileSystem.Exists(path)
}

// IsFile returns true, if the given path is a file.
func IsFile(path string) (bool, errors.Error) {
	return DefaultFileSystem.IsFile(path)
}

// IsDir returns true, if the given path is a directory.
func IsDir(path string) (bool, errors.Error) {
	return DefaultFileSystem.IsDir(path)
}

// ReadDir returns all files and directories contained in a directory.
func ReadDir(path string) ([]FileInfo, errors.Error) {
	return DefaultFileSystem.ReadDir(path)
}

// Open opens a file instance for reading and returns the handle.
func Open(path string) (File, errors.Error) {
	return DefaultFileSystem.Open(path)
}

// OpenFile opens a general purpose file instance based on flags and returns the handle.
func OpenFile(path string, flags OpenFlags) (File, errors.Error) {
	return DefaultFileSystem.OpenFile(path, flags)
}

// ReadBytes returns all bytes contained in a file.
func ReadBytes(path string) ([]byte, errors.Error) {
	return DefaultFileSystem.ReadBytes(path)
}

// ReadString returns the file content as string.
func ReadString(path string) (string, errors.Error) {
	return DefaultFileSystem.ReadString(path)
}

// ReadLines returns all lines separated by "\n", "\r" or "\r\n" from a file.
func ReadLines(path string) ([]string, errors.Error) {
	return DefaultFileSystem.ReadLines(path)
}

// CreateFile a new file (or truncate an existing) and return the file instance handle.
func CreateFile(path string) (File, errors.Error) {
	return DefaultFileSystem.CreateFile(path)
}

// CreateDirectory creates a new directory and all parent directories if they do not exist.
func CreateDirectory(path string) errors.Error {
	return DefaultFileSystem.CreateDirectory(path)
}

// WriteBytes writes all bytes to a file.
func WriteBytes(path string, content []byte) errors.Error {
	return DefaultFileSystem.WriteBytes(path, content)
}

// WriteString writes a string to a file.
func WriteString(path, content string) errors.Error {
	return DefaultFileSystem.WriteString(path, content)
}

// WriteLines writes all lines separated by the default line separator to a file.
func WriteLines(path string, lines []string) errors.Error {
	return DefaultFileSystem.WriteLines(path, lines)
}

// DeleteFile deletes a file.
func DeleteFile(path string) errors.Error {
	return DefaultFileSystem.DeleteFile(path)
}

// DeleteDirectory deletes an empty directory. If recursive is set, all contained items will be deleted first.
func DeleteDirectory(path string, recursive bool) errors.Error {
	return DefaultFileSystem.DeleteDirectory(path, recursive)
}

// MoveFile moves a file to a new location.
func MoveFile(src, dst string) errors.Error {
	return DefaultFileSystem.MoveFile(src, dst)
}

// MoveDir moves a directory to a new location.
func MoveDir(src, dst string) errors.Error {
	return DefaultFileSystem.MoveDir(src, dst)
}

// MoveAll moves all files and directories contained in src to dst.
func MoveAll(src, dst string) errors.Error {
	return DefaultFileSystem.MoveAll(src, dst)
}

// Copy clone a file or directory to the target. If the target already exists, it must be the same element type (file or directory) to be overwritten.
func Copy(src, dst string) errors.Error {
	return DefaultFileSystem.Copy(src, dst)
}

// CopyFile clones a file and overwrites the existing one.
func CopyFile(src, dst string) errors.Error {
	return DefaultFileSystem.CopyFile(src, dst)
}

// CopyDir recursively clones a directory overwriting all existing files.
func CopyDir(src, dst string) errors.Error {
	return DefaultFileSystem.CopyDir(src, dst)
}

// CopyAll copies all files and directories contained in src to dst.
func CopyAll(src, dst string) errors.Error {
	return DefaultFileSystem.CopyAll(src, dst)
}

// CleanDir removes all files and directories from a directory.
func CleanDir(path string) errors.Error {
	return DefaultFileSystem.CleanDir(path)
}

// GetTempFile returns the path to an empty temporary file.
func GetTempFile(pattern string) (string, errors.Error) {
	return DefaultFileSystem.GetTempFile(pattern)
}

// GetTempDir returns the path to an empty temporary dir.
func GetTempDir(prefix string) (string, errors.Error) {
	return DefaultFileSystem.GetTempDir(prefix)
}

// WithTempFile creates a temporary file and deletes it when f returns.
func WithTempFile(pattern string, f func(tmpFile string) errors.Error) errors.Error {
	return DefaultFileSystem.WithTempFile(pattern, f)
}

// WithTempDir creates a temporary directory and deletes it when f returns.
func WithTempDir(prefix string, f func(tmpDir string) errors.Error) errors.Error {
	return DefaultFileSystem.WithTempDir(prefix, f)
}
