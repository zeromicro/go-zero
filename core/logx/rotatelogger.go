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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/fs"
	"github.com/zeromicro/go-zero/core/lang"
)

const (
	hoursPerDay     = 24
	bufferSize      = 100
	defaultDirMode  = 0o755
	defaultFileMode = 0o600
	gzipExt         = ".gz"
	megaBytes       = 1 << 20

	// defaultCreateMode is the mode used when the log file or a backup is
	// created, the same as os.Create does.
	defaultCreateMode = 0o666

	// The log file can be written by multiple processes, the in-memory size
	// counter only knows the bytes written by this process. The real size is
	// refreshed from the file system, but at most once per sizeCheckInterval:
	// one stat per entry would be too costly, especially on a network
	// filesystem.
	sizeCheckInterval = time.Second

	// A rotation that keeps failing, for instance with no space left or on a
	// read-only backup directory, copies the whole log file on every attempt.
	// Retry at most once per rotateRetryInterval.
	rotateRetryInterval = time.Minute

	// defaultFileTimeFormat is the default time layout used in backup file
	// names. It used to be time.RFC3339, whose colons are not valid in file
	// names on Windows. The layout is colon-free now, and the nanoseconds
	// keep the backups of multiple processes rotating at the same moment
	// apart.
	defaultFileTimeFormat = "2006-01-02T15-04-05.000000000"
)

var (
	// ErrLogFileClosed is an error that indicates the log file is already closed.
	ErrLogFileClosed = errors.New("error: log file closed")

	// fileTimeFormat is the time layout used in backup file names, it can be
	// overridden with the FileTimeFormat config option.
	fileTimeFormat = defaultFileTimeFormat
)

type (
	// A RotateRule interface is used to define the log rotating rules.
	RotateRule interface {
		BackupFileName() string
		MarkRotated()
		OutdatedFiles() []string
		ShallRotate(size int64) bool
	}

	// A RotateLogger is a Logger that can rotate log files with given rules.
	RotateLogger struct {
		filename string
		backup   string
		fp       *os.File
		channel  chan []byte
		done     chan lang.PlaceholderType
		rule     RotateRule
		compress bool
		// can't use threading.RoutineGroup because of cycle import
		waitGroup         sync.WaitGroup
		closeOnce         sync.Once
		currentSize       int64
		lastSizeCheckTime time.Time
		lastRotateFail    time.Time
	}

	// A DailyRotateRule is a rule to daily rotate the log files.
	DailyRotateRule struct {
		rotatedTime string
		filename    string
		delimiter   string
		days        int
		gzip        bool
	}

	// SizeLimitRotateRule a rotation rule that makes the log file rotated based on size
	SizeLimitRotateRule struct {
		DailyRotateRule
		maxSize    int64
		maxBackups int
	}
)

// DefaultRotateRule is a default log rotating rule, currently DailyRotateRule.
func DefaultRotateRule(filename, delimiter string, days int, gzip bool) RotateRule {
	return &DailyRotateRule{
		rotatedTime: getNowDate(),
		filename:    filename,
		delimiter:   delimiter,
		days:        days,
		gzip:        gzip,
	}
}

// BackupFileName returns the backup filename on rotating.
func (r *DailyRotateRule) BackupFileName() string {
	return fmt.Sprintf("%s%s%s", r.filename, r.delimiter, getNowDate())
}

// MarkRotated marks the rotated time of r to be the current time.
func (r *DailyRotateRule) MarkRotated() {
	r.rotatedTime = getNowDate()
}

// OutdatedFiles returns the files that exceeded the keeping days.
func (r *DailyRotateRule) OutdatedFiles() []string {
	if r.days <= 0 {
		return nil
	}

	var pattern string
	if r.gzip {
		pattern = fmt.Sprintf("%s%s*%s", r.filename, r.delimiter, gzipExt)
	} else {
		pattern = fmt.Sprintf("%s%s*", r.filename, r.delimiter)
	}

	files, err := filepath.Glob(pattern)
	if err != nil {
		Errorf("failed to delete outdated log files, error: %s", err)
		return nil
	}

	var buf strings.Builder
	boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay*r.days)).Format(time.DateOnly)
	buf.WriteString(r.filename)
	buf.WriteString(r.delimiter)
	buf.WriteString(boundary)
	if r.gzip {
		buf.WriteString(gzipExt)
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

// ShallRotate checks if the file should be rotated.
func (r *DailyRotateRule) ShallRotate(_ int64) bool {
	return len(r.rotatedTime) > 0 && getNowDate() != r.rotatedTime
}

// NewSizeLimitRotateRule returns the rotation rule with size limit
func NewSizeLimitRotateRule(filename, delimiter string, days, maxSize, maxBackups int, gzip bool) RotateRule {
	return &SizeLimitRotateRule{
		DailyRotateRule: DailyRotateRule{
			rotatedTime: getNowDateInFileTimeFormat(),
			filename:    filename,
			delimiter:   delimiter,
			days:        days,
			gzip:        gzip,
		},
		maxSize:    int64(maxSize) * megaBytes,
		maxBackups: maxBackups,
	}
}

func (r *SizeLimitRotateRule) BackupFileName() string {
	dir := filepath.Dir(r.filename)
	prefix, ext := r.parseFilename()
	timestamp := getNowDateInFileTimeFormat()
	return filepath.Join(dir, fmt.Sprintf("%s%s%s%s", prefix, r.delimiter, timestamp, ext))
}

func (r *SizeLimitRotateRule) MarkRotated() {
	r.rotatedTime = getNowDateInFileTimeFormat()
}

func (r *SizeLimitRotateRule) OutdatedFiles() []string {
	dir := filepath.Dir(r.filename)
	prefix, ext := r.parseFilename()

	var pattern string
	if r.gzip {
		pattern = fmt.Sprintf("%s%s%s%s*%s%s", dir, string(filepath.Separator),
			prefix, r.delimiter, ext, gzipExt)
	} else {
		pattern = fmt.Sprintf("%s%s%s%s*%s", dir, string(filepath.Separator),
			prefix, r.delimiter, ext)
	}

	files, err := filepath.Glob(pattern)
	if err != nil {
		Errorf("failed to delete outdated log files, error: %s", err)
		return nil
	}

	// The backups are ordered by the time in their names. The time layout used
	// to be time.RFC3339, whose colons sort after the dash of the current
	// layout, so the names of both layouts don't order by time when compared
	// raw. backupNameOrderKey folds the colons into dashes, so that both
	// layouts order by time, and the order of the names without colons keeps
	// unchanged.
	orderKey := backupNameOrderKey
	sort.Slice(files, func(i, j int) bool {
		return orderKey(files[i]) < orderKey(files[j])
	})

	outdated := make(map[string]lang.PlaceholderType)

	// test if too many backups
	if r.maxBackups > 0 && len(files) > r.maxBackups {
		for _, f := range files[:len(files)-r.maxBackups] {
			outdated[f] = lang.Placeholder
		}
		files = files[len(files)-r.maxBackups:]
	}

	// test if any too old backups
	if r.days > 0 {
		boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay*r.days)).Format(fileTimeFormat)
		boundaryFile := filepath.Join(dir, fmt.Sprintf("%s%s%s%s", prefix, r.delimiter, boundary, ext))
		if r.gzip {
			boundaryFile += gzipExt
		}
		boundaryKey := orderKey(boundaryFile)
		for _, f := range files {
			if orderKey(f) >= boundaryKey {
				break
			}
			outdated[f] = lang.Placeholder
		}
	}

	result := make([]string, 0, len(outdated))
	for k := range outdated {
		result = append(result, k)
	}
	return result
}

func (r *SizeLimitRotateRule) ShallRotate(size int64) bool {
	return r.maxSize > 0 && r.maxSize < size
}

// backupNameOrderKey converts a backup filename into a form that orders by the
// time in the name when compared as a string, see SizeLimitRotateRule.OutdatedFiles.
func backupNameOrderKey(name string) string {
	return strings.ReplaceAll(name, ":", "-")
}

func (r *SizeLimitRotateRule) parseFilename() (prefix, ext string) {
	logName := filepath.Base(r.filename)
	ext = filepath.Ext(r.filename)
	prefix = logName[:len(logName)-len(ext)]
	return
}

// NewLogger returns a RotateLogger with given filename and rule, etc.
func NewLogger(filename string, rule RotateRule, compress bool) (*RotateLogger, error) {
	l := &RotateLogger{
		filename: filename,
		channel:  make(chan []byte, bufferSize),
		done:     make(chan lang.PlaceholderType),
		rule:     rule,
		compress: compress,
	}
	if err := l.initialize(); err != nil {
		return nil, err
	}

	l.startWorker()
	return l, nil
}

// Close closes l.
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
	}

	return l.backup
}

func (l *RotateLogger) initialize() error {
	l.backup = l.rule.BackupFileName()

	if fileInfo, err := os.Stat(l.filename); err != nil {
		basePath := path.Dir(l.filename)
		if _, err = os.Stat(basePath); err != nil {
			if err = os.MkdirAll(basePath, defaultDirMode); err != nil {
				return err
			}
		}

		if l.fp, err = os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, defaultCreateMode); err != nil {
			return err
		}
	} else {
		if l.fp, err = os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, defaultFileMode); err != nil {
			return err
		}

		l.currentSize = fileInfo.Size()
	}

	fs.CloseOnExec(l.fp)

	return nil
}

func (l *RotateLogger) maybeCompressFile(file string) {
	if !l.compress {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			ErrorStack(r)
		}
	}()

	if _, err := os.Stat(file); err != nil {
		// file doesn't exist or another error, ignore compression
		return
	}

	compressLogFile(file)
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
	backupFilename := l.getBackupFilename()
	if len(backupFilename) > 0 {
		if err := l.backupAndTruncate(backupFilename); err != nil {
			return err
		}
	}

	l.backup = l.rule.BackupFileName()
	return nil
}

// backupAndTruncate backs the log file up and truncates it in place. It used
// to be a rename: with multiple processes sharing the log file, the file
// descriptors opened by the other processes keep pointing to the renamed file,
// so their logs get lost into the backup and the log file starts from scratch
// for this process only. Truncating keeps the inode of the log file, so the
// descriptors of all the processes stay valid and keep appending to the file
// that is actually active.
func (l *RotateLogger) backupAndTruncate(backupFilename string) error {
	fileInfo, err := os.Stat(l.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// the file is gone, recreate it, otherwise the logs written
			// by this process would go to the unlinked file.
			return l.reopenFile()
		}
		// The file cannot be stated, so the rotation cannot be performed.
		// Report the error instead of reporting a fake success: the caller
		// must not mark the rotation done, otherwise the file grows
		// without end and nothing retries.
		return err
	}
	if fileInfo.Size() == 0 {
		// nothing to back up, another process may have rotated the file
		// at the same moment.
		return nil
	}

	// O_EXCL: two processes can rotate the file at the same moment, the one
	// that comes second must not overwrite the backup of the first one.
	dst, err := os.OpenFile(backupFilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileInfo.Mode().Perm())
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			// another process rotated the file at the same moment, it has
			// already backed up and truncated the content.
			return nil
		}
		return err
	}

	src, err := os.Open(l.filename)
	if err != nil {
		_ = dst.Close()
		_ = os.Remove(backupFilename)
		return err
	}
	// Copy the snapshot size only: on a shared volume the file keeps
	// growing while copying, an unbounded copy could run for a long time.
	// The entries appended after the snapshot get truncated anyway, this is
	// the copytruncate trade-off.
	written, err := io.CopyN(dst, src, fileInfo.Size())
	_ = src.Close()
	if errors.Is(err, io.EOF) {
		// The file was truncated by another process while copying, the
		// partial backup is kept when it is not empty.
		err = nil
	}
	// Close the backup unconditionally, a file left open cannot be removed
	// on Windows, and a stale backup would block the next attempt with
	// O_EXCL. Keep the first error.
	if closeErr := dst.Close(); err == nil && closeErr != nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(backupFilename)
		return err
	}
	if written == 0 {
		// an empty backup would take one of the slots kept by the
		// maxBackups rule.
		_ = os.Remove(backupFilename)
		return nil
	}

	if err = os.Truncate(l.filename, 0); err != nil {
		// The truncation may or may not have landed. The file being still
		// non-empty means it didn't: the backup is a copy of a file that
		// is going to be copied again, remove it, so the next attempt is
		// not blocked by the O_EXCL above and reports a fake success.
		if fi, statErr := os.Stat(l.filename); statErr == nil && fi.Size() > 0 {
			_ = os.Remove(backupFilename)
			return err
		}
		// The file is empty, the backup holds the logs, the rotation is
		// done despite the error.
		l.postRotate(backupFilename)
		return nil
	}

	l.postRotate(backupFilename)
	return nil
}

func (l *RotateLogger) reopenFile() error {
	fp, err := os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, defaultCreateMode)
	if err != nil {
		return err
	}
	// Close the descriptor of the removed file only after the new one is
	// open, a failure keeps the logger usable with the old descriptor
	// instead of leaving a nil one behind.
	if l.fp != nil {
		_ = l.fp.Close()
	}
	fs.CloseOnExec(fp)
	l.fp = fp
	l.currentSize = 0
	return nil
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
				// avoid losing logs before closing.
				for {
					select {
					case event := <-l.channel:
						l.write(event)
					default:
						return
					}
				}
			}
		}
	}()
}

func (l *RotateLogger) write(v []byte) {
	if time.Since(l.lastSizeCheckTime) >= sizeCheckInterval {
		l.lastSizeCheckTime = time.Now()
		// the file can be written by multiple processes, the in-memory
		// counter only knows the bytes written by this process, refresh
		// the real size from the file system.
		if fileInfo, err := os.Stat(l.filename); err == nil {
			l.currentSize = fileInfo.Size()
		}
	}

	if l.rule.ShallRotate(l.currentSize + int64(len(v))) {
		if time.Since(l.lastRotateFail) >= rotateRetryInterval {
			if err := l.rotate(); err != nil {
				log.Println(err)
				l.lastRotateFail = time.Now()
			} else {
				l.rule.MarkRotated()
				l.currentSize = 0
			}
		}
	}
	if l.fp != nil {
		l.fp.Write(v)
		l.currentSize += int64(len(v))
	}
}

func compressLogFile(file string) {
	start := time.Now()
	Infof("compressing log file: %s", file)
	if err := gzipFile(file, fileSys); err != nil {
		Errorf("compress error: %s", err)
	} else {
		Infof("compressed log file: %s, took %s", file, time.Since(start))
	}
}

func getNowDate() string {
	return time.Now().Format(time.DateOnly)
}

func getNowDateInFileTimeFormat() string {
	return time.Now().Format(fileTimeFormat)
}

func gzipFile(file string, fsys fileSystem) (err error) {
	in, err := fsys.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		if e := fsys.Close(in); e != nil {
			Errorf("failed to close file: %s, error: %v", file, e)
		}
		if err == nil {
			// only remove the original file when compression is successful
			err = fsys.Remove(file)
		}
	}()

	out, err := fsys.Create(fmt.Sprintf("%s%s", file, gzipExt))
	if err != nil {
		return err
	}
	defer func() {
		e := fsys.Close(out)
		if err == nil {
			err = e
		}
	}()

	w := gzip.NewWriter(out)
	if _, err = fsys.Copy(w, in); err != nil {
		// failed to copy, no need to close w
		return err
	}

	return fsys.Close(w)
}
