package filex

import (
	"io"
	"os"
)

const bufSize = 1024

func FirstLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return firstLine(file)
}

func LastLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return lastLine(filename, file)
}

func firstLine(file *os.File) (string, error) {
	var first []byte
	var offset int64
	for {
		buf := make([]byte, bufSize)
		n, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return "", err
		}

		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				return string(append(first, buf[:i]...)), nil
			}
		}

		first = append(first, buf[:n]...)
		offset += bufSize
	}
}

func lastLine(filename string, file *os.File) (string, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return "", err
	}

	var last []byte
	offset := info.Size()
	for {
		offset -= bufSize
		if offset < 0 {
			offset = 0
		}
		buf := make([]byte, bufSize)
		n, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return "", err
		}

		if buf[n-1] == '\n' {
			buf = buf[:n-1]
			n -= 1
		} else {
			buf = buf[:n]
		}
		for n -= 1; n >= 0; n-- {
			if buf[n] == '\n' {
				return string(append(buf[n+1:], last...)), nil
			}
		}

		last = append(buf, last...)
	}
}
