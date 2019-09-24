package log

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/dairaga/config"
)

// Severity ...
type Severity int

// Severity values, see https://godoc.org/google.golang.org/genproto/googleapis/logging/type#LogSeverity.
const (
	DEFAULT Severity = 100 * iota
	DEBUG
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
	ALERT
	EMERGENCY
)

func (s Severity) String() string {
	switch s {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case NOTICE:
		return "NOTICE"
	case WARNING:
		return "WARN"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	case ALERT:
		return "ALERT"
	case EMERGENCY:
		return "EMERGENCY"
	default:
		return "TRACE"
	}
}

func toSeverity(s string) Severity {
	switch strings.ToUpper(s) {
	case "TRACE":
		return DEFAULT
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "NOTICE":
		return NOTICE
	case "WARNING", "WARN":
		return WARNING
	case "ERROR":
		return ERROR
	case "CRITICAL", "FATAL":
		return CRITICAL
	case "ALERT":
		return ALERT
	case "EMERGENCY":
		return EMERGENCY
	default:
		return DEFAULT
	}
}

type logmsg struct {
	now      time.Time
	severity Severity
	file     string
	line     int
	data     interface{}
}

var (
	pipe    = make(chan logmsg, 64)
	loggers []Logger
	lock    = &sync.Mutex{}
)

func shortFile(file string) string {
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			return file[i+1:]
		}
	}
	return file
}

func caller(skip int) (string, int) {
	_, file, no, ok := runtime.Caller(skip)
	if !ok {
		return "???", 0
	}

	return shortFile(file), no
}

func output(severity Severity, data interface{}) {
	file, no := caller(3)

	pipe <- logmsg{
		now:      time.Now(),
		severity: severity,
		file:     file,
		line:     no,
		data:     data,
	}
}

// ----------------------------------------------------------------------------

// Trace ...
func Trace(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(DEFAULT, msg)
}

// Trancef ...
func Trancef(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(DEFAULT, msg)
}

// TraceJSON ...
func TraceJSON(a json.RawMessage) {
	output(DEFAULT, a)
}

// Debug ...
func Debug(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(DEBUG, msg)
}

// Debugf ...
func Debugf(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(DEBUG, msg)
}

// DebugJSON ...
func DebugJSON(a json.RawMessage) {
	output(DEBUG, a)
}

// Info ...
func Info(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(INFO, msg)
}

// Infof ...
func Infof(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(INFO, msg)
}

// Warn ...
func Warn(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(WARNING, msg)
}

// Warnf ...
func Warnf(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(WARNING, msg)
}

// Error ...
func Error(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(ERROR, msg)
}

// Errorf ...
func Errorf(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(ERROR, msg)
}

// Fatal ...
func Fatal(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(CRITICAL, msg)
}

// Fatalf ...
func Fatalf(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(CRITICAL, msg)
}

// Alert ...
func Alert(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(ALERT, msg)
}

// Alertf ...
func Alertf(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(ALERT, msg)
}

// Emergency ...
func Emergency(a ...interface{}) {
	msg := fmt.Sprint(a...)
	output(EMERGENCY, msg)
}

// Emergencyf ...
func Emergencyf(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	output(EMERGENCY, msg)
}

// ----------------------------------------------------------------------------

// Register ...
func Register(l Logger) {
	lock.Lock()
	defer lock.Unlock()
	loggers = append(loggers, l)
}

// Start start to log.
func start() {

	for {
		msg, ok := <-pipe
		if ok {
			switch v := msg.data.(type) {
			case string:
				for _, l := range loggers {
					l.Output(msg.now, msg.severity, msg.file, msg.line, v)
				}
			case json.RawMessage:
				for _, l := range loggers {
					l.OutputJSON(msg.now, msg.severity, msg.file, msg.line, v)
				}
			default:
				for _, l := range loggers {
					l.OutputStruct(msg.now, msg.severity, msg.file, msg.line, v)
				}
			}
		}
	}
}

// Close close all loggers
func Close(wait time.Duration) {
	if wait > 0 {
		time.Sleep(wait)
	}

	for _, l := range loggers {
		l.Close()
	}
	close(pipe)

}

// ----------------------------------------------------------------------------

func init() {
	tmp := config.GetString("log.root.level")

	loggers = append(loggers, &stdlogger{
		severity: toSeverity(tmp),
		out:      os.Stderr,
	})

	initGCPLogger()
	initRollingLogger()
	go start()
}
