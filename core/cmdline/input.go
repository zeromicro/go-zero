package cmdline

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnterToContinue let stdin waiting for an enter key to continue.
func EnterToContinue() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// ReadLine shows prompt to stdout and read a line from stdin.
func ReadLine(prompt string) string {
	fmt.Print(prompt)
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(input)
}
