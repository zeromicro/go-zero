package cmd

import (
	"fmt"
	"testing"
)

func Test_supportGoStdFlag(t *testing.T) {
	result := supportGoStdFlag([]string{"goctl", "api", "go", "-api", "demo.api", "-dir", "."})
	fmt.Println(result)
}
