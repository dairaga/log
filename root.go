package log

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type stdlogger struct {
	severity Severity
	out      io.WriteCloser
	mutex    *sync.Mutex
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func formatHeader(buf *[]byte, t time.Time, file string, line int) {

	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '-')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '-')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)

	*buf = append(*buf, '.')
	itoa(buf, t.Nanosecond()/1e3, 6)

	*buf = append(*buf, ' ')

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	*buf = append(*buf, file...)
	*buf = append(*buf, ':')
	itoa(buf, line, -1)
	*buf = append(*buf, ": "...)
}

func (l *stdlogger) outputstring(now time.Time, severity Severity, file string, lineNo int, message string) {
	//now = now.UTC()
	buf := make([]byte, 0, 128)
	formatHeader(&buf, now, file, lineNo)
	buf = append(buf, '[')
	buf = append(buf, severity.String()...)
	buf = append(buf, "] "...)
	buf = append(buf, message...)
	if len(buf) == 0 || buf[len(buf)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.out.Write(buf)
}

func (l *stdlogger) Output(now time.Time, severity Severity, file string, lineNo int, message string) {
	if severity >= l.severity {
		l.outputstring(now, severity, file, lineNo, message)
	}
}

func (l *stdlogger) OutputJSON(now time.Time, severity Severity, file string, lineNo int, message json.RawMessage) {
	l.Output(now, severity, file, lineNo, string(message))
}

func (l *stdlogger) OutputStruct(now time.Time, severity Severity, file string, lineNo int, message interface{}) {
	if severity >= l.severity {
		databytes, err := json.Marshal(message)
		if err != nil {
			return
		}

		l.outputstring(now, severity, file, lineNo, string(databytes))
	}
}

func (l *stdlogger) Close() error {
	return l.out.Close()
}
