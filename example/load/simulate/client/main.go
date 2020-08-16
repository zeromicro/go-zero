package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/fx"
	"github.com/tal-tech/go-zero/core/logx"
)

var (
	errServiceUnavailable = errors.New("service unavailable")
	total                 int64
	pass                  int64
	fail                  int64
	drop                  int64
	seconds               int64 = 1
)

func main() {
	flag.Parse()

	fp, err := os.Create("result.csv")
	logx.Must(err)
	defer fp.Close()
	fmt.Fprintln(fp, "seconds,total,pass,fail,drop")

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			reset(fp)
		}
	}()

	for i := 0; ; i++ {
		it := time.NewTicker(time.Second / time.Duration(atomic.LoadInt64(&seconds)))
		func() {
			for j := 0; j < int(seconds); j++ {
				<-it.C
				go issueRequest()
			}
		}()
		it.Stop()

		cur := atomic.AddInt64(&seconds, 1)
		fmt.Println(cur)
	}
}

func issueRequest() {
	atomic.AddInt64(&total, 1)
	err := fx.DoWithTimeout(func() error {
		return job()
	}, time.Second)
	switch err {
	case nil:
		atomic.AddInt64(&pass, 1)
	case errServiceUnavailable:
		atomic.AddInt64(&drop, 1)
	default:
		atomic.AddInt64(&fail, 1)
	}
}

func job() error {
	resp, err := http.Get("http://localhost:3333/")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return errServiceUnavailable
	}
}

func reset(writer io.Writer) {
	fmt.Fprintf(writer, "%d,%d,%d,%d,%d\n",
		atomic.LoadInt64(&seconds),
		atomic.SwapInt64(&total, 0),
		atomic.SwapInt64(&pass, 0),
		atomic.SwapInt64(&fail, 0),
		atomic.SwapInt64(&drop, 0),
	)
}
