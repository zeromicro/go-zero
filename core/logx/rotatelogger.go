package logx

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/fs"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/timex"
)

const (
	dateFormat      = "2006-01-02"
	hoursPerDay     = 24
	bufferSize      = 100
	defaultDirMode  = 0755
	defaultFileMode = 0600
)

var ErrLogFileClosed = errors.New("error: log file closed")

type (
	RotateRule interface {
		BackupFileName() string
		MarkRotated()
		OutdatedFiles() []string
		ShallRotate() bool
	}

	RotateLogger struct {
		filename string
		backup   string
		fp       *os.File
		channel  chan []byte
		done     chan lang.PlaceholderType
		rule     RotateRule
		compress bool
		keepDays int
		// can't use threading.RoutineGroup because of cycle import
		waitGroup sync.WaitGroup
		closeOnce sync.Once
	}

	DailyRotateRule struct {
		rotatedTime string
		filename    string
		delimiter   string
		days        int
		gzip        bool
	}
)

func DefaultRotateRule(filename, delimiter string, days int, gzip bool) RotateRule {
	return &DailyRotateRule{
		rotatedTime: getNowDate(),
		filename:    filename,
		delimiter:   delimiter,
		days:        days,
		gzip:        gzip,
	}
}

func (r *DailyRotateRule) BackupFileName() string {
	return fmt.Sprintf("%s%s%s", r.filename, r.delimiter, getNowDate())
}

func (r *DailyRotateRule) MarkRotated() {
	r.rotatedTime = getNowDate()
}

func (r *DailyRotateRule) OutdatedFiles() []string {
	if r.days <= 0 {
		return nil
	}

	var pattern string
	if r.gzip {
		pattern = fmt.Sprintf("%s%s*.gz", r.filename, r.delimiter)
	} else {
		pattern = fmt.Sprintf("%s%s*", r.filename, r.delimiter)
	}

	files, err := filepath.Glob(pattern)
	if err != nil {
		Errorf("failed to delete outdated log files, error: %s", err)
		return nil
	}

	var buf strings.Builder
	boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay*r.days)).Format(dateFormat)
	fmt.Fprintf(&buf, "%s%s%s", r.filename, r.delimiter, boundary)
	if r.gzip {
		buf.WriteString(".gz")
	}
	boundaryFile := buf.String()

	var outdates []string
	for _, file := range files {
		if file < boundaryFile {
			outdates = append(outdates, file)
		}
	}

	return outdates
}

func (r *DailyRotateRule) ShallRotate() bool {
	return len(r.rotatedTime) > 0 && getNowDate() != r.rotatedTime
}

func NewLogger(filename string, rule RotateRule, compress bool) (*RotateLogger, error) {
	l := &RotateLogger{
		filename: filename,
		channel:  make(chan []byte, bufferSize),
		done:     make(chan lang.PlaceholderType),
		rule:     rule,
		compress: compress,
	}
	if err := l.init(); err != nil {
		return nil, err
	}

	l.startWorker()
	return l, nil
}

func (l *RotateLogger) Close() error {
	var err error

	l.closeOnce.Do(func() {
		close(l.done)
		l.waitGroup.Wait()

		if err = l.fp.Sync(); err != nil {
			return
		}

		err = l.fp.Close()
	})

	return err
}

func (l *RotateLogger) Write(data []byte) (int, error) {
	select {
	case l.channel <- data:
		return len(data), nil
	case <-l.done:
		log.Println(string(data))
		return 0, ErrLogFileClosed
	}
}

func (l *RotateLogger) getBackupFilename() string {
	if len(l.backup) == 0 {
		return l.rule.BackupFileName()
	} else {
		return l.backup
	}
}

func (l *RotateLogger) init() error {
	l.backup = l.rule.BackupFileName()

	if _, err := os.Stat(l.filename); err != nil {
		basePath := path.Dir(l.filename)
		if _, err = os.Stat(basePath); err != nil {
			if err = os.MkdirAll(basePath, defaultDirMode); err != nil {
				return err
			}
		}

		if l.fp, err = os.Create(l.filename); err != nil {
			return err
		}
	} else if l.fp, err = os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, defaultFileMode); err != nil {
		return err
	}

	fs.CloseOnExec(l.fp)

	return nil
}

func (l *RotateLogger) maybeCompressFile(file string) {
	if l.compress {
		defer func() {
			if r := recover(); r != nil {
				ErrorStack(r)
			}
		}()
		compressLogFile(file)
	}
}

func (l *RotateLogger) maybeDeleteOutdatedFiles() {
	files := l.rule.OutdatedFiles()
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			Errorf("failed to remove outdated file: %s", file)
		}
	}
}

func (l *RotateLogger) postRotate(file string) {
	go func() {
		// we cannot use threading.GoSafe here, because of import cycle.
		l.maybeCompressFile(file)
		l.maybeDeleteOutdatedFiles()
	}()
}

func (l *RotateLogger) rotate() error {
	if l.fp != nil {
		err := l.fp.Close()
		l.fp = nil
		if err != nil {
			return err
		}
	}

	_, err := os.Stat(l.filename)
	if err == nil && len(l.backup) > 0 {
		backupFilename := l.getBackupFilename()
		err = os.Rename(l.filename, backupFilename)
		if err != nil {
			return err
		}

		l.postRotate(backupFilename)
	}

	l.backup = l.rule.BackupFileName()
	if l.fp, err = os.Create(l.filename); err == nil {
		fs.CloseOnExec(l.fp)
	}

	return err
}

func (l *RotateLogger) startWorker() {
	l.waitGroup.Add(1)

	go func() {
		defer l.waitGroup.Done()

		for {
			select {
			case event := <-l.channel:
				l.write(event)
			case <-l.done:
				return
			}
		}
	}()
}

func (l *RotateLogger) write(v []byte) {
	if l.rule.ShallRotate() {
		if err := l.rotate(); err != nil {
			log.Println(err)
		} else {
			l.rule.MarkRotated()
		}
	}
	if l.fp != nil {
		l.fp.Write(v)
	}
}

func compressLogFile(file string) {
	start := timex.Now()
	Infof("compressing log file: %s", file)
	if err := gzipFile(file); err != nil {
		Errorf("compress error: %s", err)
	} else {
		Infof("compressed log file: %s, took %s", file, timex.Since(start))
	}
}

func getNowDate() string {
	return time.Now().Format(dateFormat)
}

func gzipFile(file string) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(fmt.Sprintf("%s.gz", file))
	if err != nil {
		return err
	}
	defer out.Close()

	w := gzip.NewWriter(out)
	if _, err = io.Copy(w, in); err != nil {
		return err
	} else if err = w.Close(); err != nil {
		return err
	}

	return os.Remove(file)
}
