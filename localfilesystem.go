package fs

// LocalFileSystem allows access to the file system of the host machine.
type LocalFileSystem struct {
	root string
}

// NewLocalRootFileSystem returns a new local file system at root.
func NewLocalRootFileSystem() *LocalFileSystem {
	return NewLocalFileSystem("")
}

// NewLocalFileSystem returns a new local file system rooted at the given directory. Access to parent directories is prohibited.
func NewLocalFileSystem(root string) *LocalFileSystem {
	return &LocalFileSystem{root}
}
