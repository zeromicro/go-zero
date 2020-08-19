package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/threading"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	freq     = flag.Int("freq", 100, "frequency")
	duration = flag.String("duration", "10s", "duration")
)

type (
	counting struct {
		ok      int
		fail    int
		reject  int
		errs    int
		unknown int
	}

	metric struct {
		counting
		lock sync.Mutex
	}
)

func (m *metric) addOk() {
	m.lock.Lock()
	m.ok++
	m.lock.Unlock()
}

func (m *metric) addFail() {
	m.lock.Lock()
	m.ok++
	m.lock.Unlock()
}

func (m *metric) addReject() {
	m.lock.Lock()
	m.ok++
	m.lock.Unlock()
}

func (m *metric) addErrs() {
	m.lock.Lock()
	m.errs++
	m.lock.Unlock()
}

func (m *metric) addUnknown() {
	m.lock.Lock()
	m.unknown++
	m.lock.Unlock()
}

func (m *metric) reset() counting {
	m.lock.Lock()
	result := counting{
		ok:      m.ok,
		fail:    m.fail,
		reject:  m.reject,
		errs:    m.errs,
		unknown: m.unknown,
	}

	m.ok = 0
	m.fail = 0
	m.reject = 0
	m.errs = 0
	m.unknown = 0
	m.lock.Unlock()

	return result
}

func runRequests(url string, frequency int, metrics *metric, done <-chan lang.PlaceholderType) {
	ticker := time.NewTicker(time.Second / time.Duration(frequency))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				resp, err := http.Get(url)
				if err != nil {
					metrics.addErrs()
					return
				}
				defer resp.Body.Close()

				switch resp.StatusCode {
				case http.StatusOK:
					metrics.addOk()
				case http.StatusInternalServerError:
					metrics.addFail()
				case http.StatusServiceUnavailable:
					metrics.addReject()
				default:
					metrics.addUnknown()
				}
			}()
		case <-done:
			return
		}
	}
}

func main() {
	flag.Parse()

	fp, err := os.Create("result.csv")
	logx.Must(err)
	defer fp.Close()
	fmt.Fprintln(fp, "seconds,goodOk,goodFail,goodReject,goodErrs,goodUnknowns,goodDropRatio,"+
		"heavyOk,heavyFail,heavyReject,heavyErrs,heavyUnknowns,heavyDropRatio")

	var gm, hm metric
	dur, err := time.ParseDuration(*duration)
	logx.Must(err)
	done := make(chan lang.PlaceholderType)
	group := threading.NewRoutineGroup()
	group.RunSafe(func() {
		runRequests("http://localhost:8080/heavy", *freq, &hm, done)
	})
	group.RunSafe(func() {
		runRequests("http://localhost:8080/good", *freq, &gm, done)
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		var seconds int
		for range ticker.C {
			seconds++
			g := gm.reset()
			h := hm.reset()
			fmt.Fprintf(fp, "%d,%d,%d,%d,%d,%d,%.1f,%d,%d,%d,%d,%d,%.1f\n",
				seconds, g.ok, g.fail, g.reject, g.errs, g.unknown,
				float32(g.reject)/float32(g.ok+g.fail+g.reject+g.unknown),
				h.ok, h.fail, h.reject, h.errs, h.unknown,
				float32(h.reject)/float32(h.ok+h.fail+h.reject+h.unknown))
		}
	}()

	go func() {
		bar := pb.New(int(dur / time.Second)).Start()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			bar.Increment()
		}
		bar.Finish()
	}()

	<-time.After(dur)
	close(done)
	group.Wait()
	time.Sleep(time.Millisecond * 900)
}
