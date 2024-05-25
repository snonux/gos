package vfs

// virtual file system - useful for testing as well
type VFS interface {
	ReadFile(name string) ([]byte, error)
	SaveFile(filePath string, bytes []byte) error
	FindFiles(dataPath, suffix string) ([]string, error)
}
