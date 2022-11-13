package format

import (
	_ "embed"
	"log"
	"os"
	"testing"
)

//go:embed testdata/test_format.api
var testFormatInput []byte

func TestFormat(t *testing.T) {
	t.Run("format", func(t *testing.T) { // EXPERIMENTAL: just for testing output formatter.
		err := Format(testFormatInput, os.Stdout)
		if err != nil {
			log.Fatalln(err)
		}
	})
}
