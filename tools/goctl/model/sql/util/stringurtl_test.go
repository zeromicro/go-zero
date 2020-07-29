package util

import (
	"fmt"
	"testing"
)

func TestFormatField(t *testing.T) {
	var in = "go_java"
	snakeCase, upperCase, lowerCase := FormatField(in)
	fmt.Println(snakeCase, upperCase, lowerCase)
}
