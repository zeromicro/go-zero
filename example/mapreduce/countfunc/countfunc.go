package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/gops/agent"
	"github.com/tal-tech/go-zero/core/mr"
)

var (
	dir        = flag.String("d", "", "dir to enumerate")
	stopOnFile = flag.String("s", "", "stop when got file")
	maxFiles   = flag.Int("m", 0, "at most files to process")
	mode       = flag.String("mode", "", "simulate mode, can be return|panic")
	count      uint32
)

func enumerateLines(filename string) chan string {
	output := make(chan string)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}

			if !strings.HasPrefix(line, "#") {
				output <- line
			}
		}
		close(output)
	}()
	return output
}

func mapper(filename interface{}, writer mr.Writer, cancel func(error)) {
	if len(*stopOnFile) > 0 && path.Base(filename.(string)) == *stopOnFile {
		fmt.Printf("Stop on file: %s\n", *stopOnFile)
		cancel(errors.New("stop on file"))
		return
	}

	var result int
	for line := range enumerateLines(filename.(string)) {
		if strings.HasPrefix(strings.TrimSpace(line), "func") {
			result++
		}
	}

	switch *mode {
	case "return":
		if atomic.AddUint32(&count, 1)%10 == 0 {
			return
		}
	case "panic":
		if atomic.AddUint32(&count, 1)%10 == 0 {
			panic("wow")
		}
	}

	writer.Write(result)
}

func reducer(input <-chan interface{}, writer mr.Writer, cancel func(error)) {
	var result int

	for count := range input {
		v := count.(int)
		if *maxFiles > 0 && result >= *maxFiles {
			fmt.Printf("Reached max files: %d\n", *maxFiles)
			cancel(errors.New("max files reached"))
			return
		}
		result += v
	}

	writer.Write(result)
}

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	if len(*dir) == 0 {
		flag.Usage()
	}

	fmt.Println("Processing, please wait...")

	start := time.Now()
	result, err := mr.MapReduce(func(source chan<- interface{}) {
		filepath.Walk(*dir, func(fpath string, f os.FileInfo, err error) error {
			if !f.IsDir() && path.Ext(fpath) == ".go" {
				source <- fpath
			}
			return nil
		})
	}, mapper, reducer)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
		fmt.Println("Elapsed:", time.Since(start))
		fmt.Println("Done")
	}
}
