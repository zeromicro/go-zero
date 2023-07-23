package logx

import (
	"io"
	"os"
)

var fileSys realFileSystem

type (
	fileSystem interface {
		Close(closer io.Closer) error
		Copy(writer io.Writer, reader io.Reader) (int64, error)
		Create(name string) (*os.File, error)
		Open(name string) (*os.File, error)
		Remove(name string) error
	}

	realFileSystem struct{}
)

func (fs realFileSystem) Close(closer io.Closer) error {
	return closer.Close()
}

func (fs realFileSystem) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	return io.Copy(writer, reader)
}

func (fs realFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fs realFileSystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (fs realFileSystem) Remove(name string) error {
	return os.Remove(name)
}
