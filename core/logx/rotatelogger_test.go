package logx

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestDailyRotateRuleMarkRotated(t *testing.T) {
	t.Run("daily rule", func(t *testing.T) {
		var rule DailyRotateRule
		rule.MarkRotated()
		assert.Equal(t, getNowDate(), rule.rotatedTime)
	})

	t.Run("daily rule", func(t *testing.T) {
		rule := DefaultRotateRule("test", "-", 1, false)
		_, ok := rule.(*DailyRotateRule)
		assert.True(t, ok)
	})
}

func TestDailyRotateRuleOutdatedFiles(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		var rule DailyRotateRule
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("bad files", func(t *testing.T) {
		rule := DailyRotateRule{
			filename: "[a-z",
		}
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("temp files", func(t *testing.T) {
		boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay) * 2).Format(dateFormat)
		f1, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary)
		assert.NoError(t, err)
		_ = f1.Close()
		f2, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary)
		assert.NoError(t, err)
		_ = f2.Close()
		t.Cleanup(func() {
			_ = os.Remove(f1.Name())
			_ = os.Remove(f2.Name())
		})
		rule := DailyRotateRule{
			filename: path.Join(os.TempDir(), "go-zero-test-"),
			days:     1,
		}
		assert.NotEmpty(t, rule.OutdatedFiles())
	})
}

func TestDailyRotateRuleShallRotate(t *testing.T) {
	var rule DailyRotateRule
	rule.rotatedTime = time.Now().Add(time.Hour * 24).Format(dateFormat)
	assert.True(t, rule.ShallRotate(0))
}

func TestSizeLimitRotateRuleMarkRotated(t *testing.T) {
	t.Run("size limit rule", func(t *testing.T) {
		var rule SizeLimitRotateRule
		rule.MarkRotated()
		assert.Equal(t, getNowDateInRFC3339Format(), rule.rotatedTime)
	})

	t.Run("size limit rule", func(t *testing.T) {
		rule := NewSizeLimitRotateRule("foo", "-", 1, 1, 1, false)
		rule.MarkRotated()
		assert.Equal(t, getNowDateInRFC3339Format(), rule.(*SizeLimitRotateRule).rotatedTime)
	})
}

func TestSizeLimitRotateRuleOutdatedFiles(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		var rule SizeLimitRotateRule
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
		rule.maxBackups = 0
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("bad files", func(t *testing.T) {
		rule := SizeLimitRotateRule{
			DailyRotateRule: DailyRotateRule{
				filename: "[a-z",
			},
		}
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("temp files", func(t *testing.T) {
		boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay) * 2).Format(dateFormat)
		f1, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary)
		assert.NoError(t, err)
		f2, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary)
		assert.NoError(t, err)
		boundary1 := time.Now().Add(time.Hour * time.Duration(hoursPerDay) * 2).Format(dateFormat)
		f3, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary1)
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = f1.Close()
			_ = os.Remove(f1.Name())
			_ = f2.Close()
			_ = os.Remove(f2.Name())
			_ = f3.Close()
			_ = os.Remove(f3.Name())
		})
		rule := SizeLimitRotateRule{
			DailyRotateRule: DailyRotateRule{
				filename: path.Join(os.TempDir(), "go-zero-test-"),
				days:     1,
			},
			maxBackups: 3,
		}
		assert.NotEmpty(t, rule.OutdatedFiles())
	})

	t.Run("no backups", func(t *testing.T) {
		boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay) * 2).Format(dateFormat)
		f1, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary)
		assert.NoError(t, err)
		f2, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary)
		assert.NoError(t, err)
		boundary1 := time.Now().Add(time.Hour * time.Duration(hoursPerDay) * 2).Format(dateFormat)
		f3, err := os.CreateTemp(os.TempDir(), "go-zero-test-"+boundary1)
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = f1.Close()
			_ = os.Remove(f1.Name())
			_ = f2.Close()
			_ = os.Remove(f2.Name())
			_ = f3.Close()
			_ = os.Remove(f3.Name())
		})
		rule := SizeLimitRotateRule{
			DailyRotateRule: DailyRotateRule{
				filename: path.Join(os.TempDir(), "go-zero-test-"),
				days:     1,
			},
		}
		assert.NotEmpty(t, rule.OutdatedFiles())

		logger := new(RotateLogger)
		logger.rule = &rule
		logger.maybeDeleteOutdatedFiles()
		assert.Empty(t, rule.OutdatedFiles())
	})
}

func TestSizeLimitRotateRuleShallRotate(t *testing.T) {
	var rule SizeLimitRotateRule
	rule.rotatedTime = time.Now().Add(time.Hour * 24).Format(fileTimeFormat)
	rule.maxSize = 0
	assert.False(t, rule.ShallRotate(0))
	rule.maxSize = 100
	assert.False(t, rule.ShallRotate(0))
	assert.True(t, rule.ShallRotate(101*megaBytes))
}

func TestRotateLoggerClose(t *testing.T) {
	t.Run("close", func(t *testing.T) {
		filename, err := fs.TempFilenameWithText("foo")
		assert.Nil(t, err)
		if len(filename) > 0 {
			defer os.Remove(filename)
		}
		logger, err := NewLogger(filename, new(DailyRotateRule), false)
		assert.Nil(t, err)
		_, err = logger.Write([]byte("foo"))
		assert.Nil(t, err)
		assert.Nil(t, logger.Close())
	})

	t.Run("close and write", func(t *testing.T) {
		logger := new(RotateLogger)
		logger.done = make(chan struct{})
		close(logger.done)
		_, err := logger.Write([]byte("foo"))
		assert.ErrorIs(t, err, ErrLogFileClosed)
	})

	t.Run("close without losing logs", func(t *testing.T) {
		text := "foo"
		filename, err := fs.TempFilenameWithText(text)
		assert.Nil(t, err)
		if len(filename) > 0 {
			defer os.Remove(filename)
		}
		logger, err := NewLogger(filename, new(DailyRotateRule), false)
		assert.Nil(t, err)
		msg := []byte("foo")
		n := 100
		for i := 0; i < n; i++ {
			_, err = logger.Write(msg)
			assert.Nil(t, err)
		}
		assert.Nil(t, logger.Close())
		bs, err := os.ReadFile(filename)
		assert.Nil(t, err)
		assert.Equal(t, len(msg)*n+len(text), len(bs))
	})
}

func TestRotateLoggerGetBackupFilename(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(DailyRotateRule), false)
	assert.Nil(t, err)
	assert.True(t, len(logger.getBackupFilename()) > 0)
	logger.backup = ""
	assert.True(t, len(logger.getBackupFilename()) > 0)
}

func TestRotateLoggerMayCompressFile(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(DailyRotateRule), false)
	assert.Nil(t, err)
	logger.maybeCompressFile(filename)
	_, err = os.Stat(filename)
	assert.Nil(t, err)
}

func TestRotateLoggerMayCompressFileTrue(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	logger, err := NewLogger(filename, new(DailyRotateRule), true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
	}
	logger.maybeCompressFile(filename)
	_, err = os.Stat(filename)
	assert.NotNil(t, err)
}

func TestRotateLoggerRotate(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	logger, err := NewLogger(filename, new(DailyRotateRule), true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer func() {
			os.Remove(logger.getBackupFilename())
			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
		}()
	}
	err = logger.rotate()
	switch v := err.(type) {
	case *os.LinkError:
		// avoid rename error on docker container
		assert.Equal(t, syscall.EXDEV, v.Err)
	case *os.PathError:
		// ignore remove error for tests,
		// files are cleaned in GitHub actions.
		assert.Equal(t, "remove", v.Op)
	default:
		assert.Nil(t, err)
	}
}

func TestRotateLoggerWrite(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	rule := new(DailyRotateRule)
	logger, err := NewLogger(filename, rule, true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer func() {
			os.Remove(logger.getBackupFilename())
			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
		}()
	}
	// the following write calls cannot be changed to Write, because of DATA RACE.
	logger.write([]byte(`foo`))
	rule.rotatedTime = time.Now().Add(-time.Hour * 24).Format(dateFormat)
	logger.write([]byte(`bar`))
	logger.Close()
	logger.write([]byte(`baz`))
}

func TestLogWriterClose(t *testing.T) {
	assert.Nil(t, newLogWriter(nil).Close())
}

func TestRotateLoggerWithSizeLimitRotateRuleClose(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), false)
	assert.Nil(t, err)
	_ = logger.Close()
}

func TestRotateLoggerGetBackupWithSizeLimitRotateRuleFilename(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), false)
	assert.Nil(t, err)
	assert.True(t, len(logger.getBackupFilename()) > 0)
	logger.backup = ""
	assert.True(t, len(logger.getBackupFilename()) > 0)
}

func TestRotateLoggerWithSizeLimitRotateRuleMayCompressFile(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), false)
	assert.Nil(t, err)
	logger.maybeCompressFile(filename)
	_, err = os.Stat(filename)
	assert.Nil(t, err)
}

func TestRotateLoggerWithSizeLimitRotateRuleMayCompressFileTrue(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
	}
	logger.maybeCompressFile(filename)
	_, err = os.Stat(filename)
	assert.NotNil(t, err)
}

func TestRotateLoggerWithSizeLimitRotateRuleMayCompressFileFailed(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename := stringx.RandId()
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), true)
	defer os.Remove(filename)
	if assert.NoError(t, err) {
		assert.NotPanics(t, func() {
			logger.maybeCompressFile(stringx.RandId())
		})
	}
}

func TestRotateLoggerWithSizeLimitRotateRuleRotate(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer func() {
			os.Remove(logger.getBackupFilename())
			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
		}()
	}
	err = logger.rotate()
	switch v := err.(type) {
	case *os.LinkError:
		// avoid rename error on docker container
		assert.Equal(t, syscall.EXDEV, v.Err)
	case *os.PathError:
		// ignore remove error for tests,
		// files are cleaned in GitHub actions.
		assert.Equal(t, "remove", v.Op)
	default:
		assert.Nil(t, err)
	}
}

func TestRotateLoggerWithSizeLimitRotateRuleWrite(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	rule := new(SizeLimitRotateRule)
	logger, err := NewLogger(filename, rule, true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer func() {
			os.Remove(logger.getBackupFilename())
			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
		}()
	}
	// the following write calls cannot be changed to Write, because of DATA RACE.
	logger.write([]byte(`foo`))
	rule.rotatedTime = time.Now().Add(-time.Hour * 24).Format(dateFormat)
	logger.write([]byte(`bar`))
	logger.Close()
	logger.write([]byte(`baz`))
}

func TestGzipFile(t *testing.T) {
	err := errors.New("any error")

	t.Run("gzip file open failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			openFn: func(name string) (*os.File, error) {
				return nil, err
			},
		}
		assert.ErrorIs(t, err, gzipFile("any", fsys))
		assert.False(t, fsys.Removed())
	})

	t.Run("gzip file create failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			createFn: func(name string) (*os.File, error) {
				return nil, err
			},
		}
		assert.ErrorIs(t, err, gzipFile("any", fsys))
		assert.False(t, fsys.Removed())
	})

	t.Run("gzip file copy failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			copyFn: func(writer io.Writer, reader io.Reader) (int64, error) {
				return 0, err
			},
		}
		assert.ErrorIs(t, err, gzipFile("any", fsys))
		assert.False(t, fsys.Removed())
	})

	t.Run("gzip file last close failed", func(t *testing.T) {
		var called int32
		fsys := &fakeFileSystem{
			closeFn: func(closer io.Closer) error {
				if atomic.AddInt32(&called, 1) > 2 {
					return err
				}
				return nil
			},
		}
		assert.NoError(t, gzipFile("any", fsys))
		assert.True(t, fsys.Removed())
	})

	t.Run("gzip file remove failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			removeFn: func(name string) error {
				return err
			},
		}
		assert.Error(t, err, gzipFile("any", fsys))
		assert.True(t, fsys.Removed())
	})

	t.Run("gzip file everything ok", func(t *testing.T) {
		fsys := &fakeFileSystem{}
		assert.NoError(t, gzipFile("any", fsys))
		assert.True(t, fsys.Removed())
	})
}

func TestRotateLogger_WithExistingFile(t *testing.T) {
	const body = "foo"
	filename, err := fs.TempFilenameWithText(body)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}

	rule := NewSizeLimitRotateRule(filename, "-", 1, 100, 3, false)
	logger, err := NewLogger(filename, rule, false)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(body)), logger.currentSize)
	assert.Nil(t, logger.Close())
}

func BenchmarkRotateLogger(b *testing.B) {
	filename := "./test.log"
	filename2 := "./test2.log"
	dailyRotateRuleLogger, err1 := NewLogger(
		filename,
		DefaultRotateRule(
			filename,
			backupFileDelimiter,
			1,
			true,
		),
		true,
	)
	if err1 != nil {
		b.Logf("Failed to new daily rotate rule logger: %v", err1)
		b.FailNow()
	}
	sizeLimitRotateRuleLogger, err2 := NewLogger(
		filename2,
		NewSizeLimitRotateRule(
			filename,
			backupFileDelimiter,
			1,
			100,
			10,
			true,
		),
		true,
	)
	if err2 != nil {
		b.Logf("Failed to new size limit rotate rule logger: %v", err1)
		b.FailNow()
	}
	defer func() {
		dailyRotateRuleLogger.Close()
		sizeLimitRotateRuleLogger.Close()
		os.Remove(filename)
		os.Remove(filename2)
	}()

	b.Run("daily rotate rule", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dailyRotateRuleLogger.write([]byte("testing\ntesting\n"))
		}
	})
	b.Run("size limit rotate rule", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sizeLimitRotateRuleLogger.write([]byte("testing\ntesting\n"))
		}
	})
}

type fakeFileSystem struct {
	removed  int32
	closeFn  func(closer io.Closer) error
	copyFn   func(writer io.Writer, reader io.Reader) (int64, error)
	createFn func(name string) (*os.File, error)
	openFn   func(name string) (*os.File, error)
	removeFn func(name string) error
}

func (f *fakeFileSystem) Close(closer io.Closer) error {
	if f.closeFn != nil {
		return f.closeFn(closer)
	}
	return nil
}

func (f *fakeFileSystem) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	if f.copyFn != nil {
		return f.copyFn(writer, reader)
	}
	return 0, nil
}

func (f *fakeFileSystem) Create(name string) (*os.File, error) {
	if f.createFn != nil {
		return f.createFn(name)
	}
	return nil, nil
}

func (f *fakeFileSystem) Open(name string) (*os.File, error) {
	if f.openFn != nil {
		return f.openFn(name)
	}
	return nil, nil
}

func (f *fakeFileSystem) Remove(name string) error {
	atomic.AddInt32(&f.removed, 1)

	if f.removeFn != nil {
		return f.removeFn(name)
	}
	return nil
}

func (f *fakeFileSystem) Removed() bool {
	return atomic.LoadInt32(&f.removed) > 0
}
