package filesys

import (
	"io"
	"os"
	"sync/atomic"
)

type fakeFileSystem struct {
	removed  int32
	closeFn  func(closer io.Closer) error
	copyFn   func(writer io.Writer, reader io.Reader) (int64, error)
	createFn func(name string) (*os.File, error)
	openFn   func(name string) (*os.File, error)
	removeFn func(name string) error
}

func (f *fakeFileSystem) Close(closer io.Closer) error {
	if f.closeFn != nil {
		return f.closeFn(closer)
	}
	return nil
}

func (f *fakeFileSystem) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	if f.copyFn != nil {
		return f.copyFn(writer, reader)
	}
	return 0, nil
}

func (f *fakeFileSystem) Create(name string) (*os.File, error) {
	if f.createFn != nil {
		return f.createFn(name)
	}
	return nil, nil
}

func (f *fakeFileSystem) Open(name string) (*os.File, error) {
	if f.openFn != nil {
		return f.openFn(name)
	}
	return nil, nil
}

func (f *fakeFileSystem) Remove(name string) error {
	atomic.AddInt32(&f.removed, 1)

	if f.removeFn != nil {
		return f.removeFn(name)
	}
	return nil
}

func (f *fakeFileSystem) Removed() bool {
	return atomic.LoadInt32(&f.removed) > 0
}
