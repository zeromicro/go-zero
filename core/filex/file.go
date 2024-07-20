package filex

import (
	"io"
	"os"
)

const bufSize = 1024

// FirstLine returns the first line of the file.
func FirstLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return firstLine(file)
}

// LastLine returns the last line of the file.
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

		if err == io.EOF {
			return string(append(first, buf[:n]...)), nil
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
	bufLen := int64(bufSize)
	offset := info.Size()

	for offset > 0 {
		if offset < bufLen {
			bufLen = offset
			offset = 0
		} else {
			offset -= bufLen
		}

		buf := make([]byte, bufLen)
		n, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return "", err
		}

		if n == 0 {
			break
		}

		if buf[n-1] == '\n' {
			buf = buf[:n-1]
			n--
		} else {
			buf = buf[:n]
		}

		for i := n - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				return string(append(buf[i+1:], last...)), nil
			}
		}

		last = append(buf, last...)
	}

	return string(last), nil
}
