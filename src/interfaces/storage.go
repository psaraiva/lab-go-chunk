package interfaces

import "os"

type ServiceStorage interface {
	CreateFile([]byte, string) error
	Clear() error
	GetFile(string) (*os.File, error)
	RemoveFile(string) error
}
