package interop

import (
	"github.com/sbreitf1/fs"

	"github.com/sbreitf1/errors"
)

// Move moves a file or directory from one file system to another recursively.
func Move(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanWrite() {
		return fs.ErrNotSupported.Msg("Source file system does not support writing").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	isFile, err := fsSrc.IsFile(src)
	if err != nil {
		return err
	}
	if isFile {
		return MoveFile(fsSrc, src, fsDst, dst)
	}

	isDir, err := fsSrc.IsDir(src)
	if err != nil {
		return err
	}
	if isDir {
		return MoveDir(fsSrc, src, fsDst, dst)
	}

	return fs.ErrNotExists.Args(src).Make()
}

// MoveFile moves a file from one file system to another.
func MoveFile(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanWrite() {
		return fs.ErrNotSupported.Msg("Source file system does not support writing").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	if err := CopyFile(fsSrc, src, fsDst, dst); err != nil {
		return err
	}

	return fsSrc.DeleteFile(src)
}

// MoveDir moves a directory recursively from one file system to another.
func MoveDir(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanWrite() {
		return fs.ErrNotSupported.Msg("Source file system does not support writing").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	if err := CopyDir(fsSrc, src, fsDst, dst); err != nil {
		return err
	}

	return fsSrc.DeleteDirectory(src, true)
}

// MoveAll moves the content of a directory to another directory recursively.
func MoveAll(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanWrite() {
		return fs.ErrNotSupported.Msg("Source file system does not support writing").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	if err := CopyAll(fsSrc, src, fsDst, dst); err != nil {
		return err
	}

	return fsSrc.CleanDir(src)
}
