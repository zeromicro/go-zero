package iox

import (
	"bytes"
	"errors"
	"io"
	"os"
)

const bufSize = 32 * 1024

// CountLines returns the number of lines in the file.
func CountLines(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var noEol bool
	buf := make([]byte, bufSize)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case errors.Is(err, io.EOF):
			if noEol {
				count++
			}
			return count, nil
		case err != nil:
			return count, err
		}

		noEol = c > 0 && buf[c-1] != '\n'
	}
}
