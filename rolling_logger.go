package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dairaga/config"
)

const (
	fileAppendMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	filePerm       = 0644
)

type rollingwriter struct {
	start  time.Time
	file   *os.File
	path   string
	prefix string
}

var _ io.WriteCloser = &rollingwriter{}

func (w *rollingwriter) reset(now time.Time) {
	w.Close()

	filename := rollingFileName(w.path, w.prefix, now)
	file, err := os.OpenFile(filename, fileAppendMode, filePerm)

	if err != nil {
		fmt.Printf("open log file %s: %v", filename, err)
		return
	}
	w.start = now
	w.file = file
}

func (w *rollingwriter) Close() error {
	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}

	return nil
}

func (w *rollingwriter) Write(p []byte) (n int, err error) {
	if w.file == nil {
		return 0, fmt.Errorf("file is nil")
	}
	now := time.Now()
	y1, m1, d1 := now.Date()
	y2, m2, d2 := w.start.Date()

	if y1 != y2 || m1 != m2 || d1 != d2 {
		w.reset(now)
	}

	return w.file.Write(p)
}

// ----------------------------------------------------------------------------

func rollingFileName(path, prefix string, start time.Time) string {
	y, m, d := start.Date()
	return fmt.Sprintf("%s/%s%d%02d%02d.log", path, prefix, y, int(m), d)
}

func initRollingLogger() {
	tmp := config.GetString("log.rolling.level")
	if tmp == "" {
		return
	}

	start := time.Now()
	path := config.GetString("log.rolling.path")
	if path != "" {
		os.MkdirAll(path, os.ModePerm)
	} else {
		path = "."
	}

	prefix := config.GetString("log.rolling.prefix", "")
	filename := rollingFileName(path, prefix, start)

	file, err := os.OpenFile(filename, fileAppendMode, filePerm)
	if err != nil {
		fmt.Printf("open log file %s: %v", filename, err)
		return
	}

	Register(&stdlogger{
		severity: toSeverity(tmp),
		out: &rollingwriter{
			start:  start,
			path:   path,
			prefix: prefix,
			file:   file,
		},
	})
}
