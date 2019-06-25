package fs

type ReadFileSystem interface {
}

type ReadWriteFileSystem interface {
	ReadFileSystem
}

type FileSystem interface {
	ReadWriteFileSystem
}
