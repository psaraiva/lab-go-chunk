package interfaces

import (
	"os"
)

type ServiceTemporaryArea interface {
	Clear() error
	CreateFileByFileSource(string, *os.File) error
	GetFile(string) (*os.File, error)
	RemoveFile(string) error
}
