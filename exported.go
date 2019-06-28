package fs

var (
	// DefaultFileSystem denots the file system used for all default accessors.
	DefaultFileSystem *FileSystem
)

func init() {
	DefaultFileSystem = NewLocal()
}
