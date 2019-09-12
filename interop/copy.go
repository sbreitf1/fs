package interop

import (
	"io"

	"github.com/sbreitf1/fs"
	"github.com/sbreitf1/fs/path"

	"github.com/sbreitf1/errors"
)

// Copy copies a file or directory from one file system to another recursively.
func Copy(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanRead() {
		return fs.ErrNotSupported.Msg("Source file system does not support reading").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	isFile, err := fsSrc.IsFile(src)
	if err != nil {
		return err
	}
	if isFile {
		return copyFile(fsSrc, src, fsDst, dst)
	}

	isDir, err := fsSrc.IsDir(src)
	if err != nil {
		return err
	}
	if isDir {
		return copyDir(fsSrc, src, fsDst, dst)
	}

	return fs.ErrNotExists.Args(src).Make()
}

// CopyFile copies a file from one file system to another.
func CopyFile(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanRead() {
		return fs.ErrNotSupported.Msg("Source file system does not support reading").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	return copyFile(fsSrc, src, fsDst, dst)
}

func copyFile(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	fSrc, err := fsSrc.Open(src)
	if err != nil {
		return err
	}
	defer fSrc.Close()

	fDst, err := fsDst.CreateFile(dst)
	if err != nil {
		return err
	}
	defer fDst.Close()

	if _, err := io.Copy(fDst, fSrc); err != nil {
		return fs.Err.Msg("Failed to copy data").Make().Cause(err)
	}

	return nil
}

// CopyDir copies a directory recursively from one file system to another.
func CopyDir(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanRead() {
		return fs.ErrNotSupported.Msg("Source file system does not support reading").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	return copyDir(fsSrc, src, fsDst, dst)
}

func copyDir(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	fsDst.CreateDirectory(dst)
	return copyAll(fsSrc, src, fsDst, dst)
}

// CopyAll copies the content of a directory to another directory recursively.
func CopyAll(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	if !fsSrc.CanRead() {
		return fs.ErrNotSupported.Msg("Source file system does not support reading").Make()
	}
	if !fsDst.CanWrite() {
		return fs.ErrNotSupported.Msg("Destination file system does not support writing").Make()
	}

	return copyAll(fsSrc, src, fsDst, dst)
}

func copyAll(fsSrc *fs.FileSystem, src string, fsDst *fs.FileSystem, dst string) errors.Error {
	files, err := fsSrc.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if err := copyDir(fsSrc, path.Join(src, f.Name()), fsDst, path.Join(dst, f.Name())); err != nil {
				return err
			}
		} else {
			if err := CopyFile(fsSrc, path.Join(src, f.Name()), fsDst, path.Join(dst, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}
