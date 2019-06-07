package log

import (
	"fmt"
	"io"
	golog "log"
	"os"
	"strings"
)

// Level ...
type Level uint8

// Level
const (
	LvAll Level = 1 + iota
	LvTrace
	LvDebug
	LvInfo
	LvWarn
	LvError
	LvFatal
)

// DefaultFlags ...
const DefaultFlags = golog.LstdFlags | golog.Lshortfile | golog.LUTC

func toLevel(x string) Level {
	switch strings.ToLower(x) {
	case "all":
		return LvAll
	case "trace":
		return LvTrace
	case "debug":
		return LvDebug
	case "info":
		return LvInfo
	case "warn":
		return LvWarn
	case "error":
		return LvError
	case "fatal", "critical":
		return LvFatal
	default:
		return LvAll
	}
}

var level = LvAll
var root *golog.Logger

func init() {
	tmp, _ := os.LookupEnv("LOG_ROOT_LEVEL")
	level = toLevel(tmp)
	root = golog.New(os.Stderr, "", DefaultFlags)
}

// SetLevel ...
func SetLevel(lv Level) {
	level = lv
}

// SetOutput ...
func SetOutput(w io.Writer) {
	root.SetOutput(w)
}

// Trace ...
func Trace(a ...interface{}) {
	if level <= LvTrace {
		root.Output(2, "[TRACE] "+fmt.Sprint(a...))
	}
}

// Tracef ...
func Tracef(f string, a ...interface{}) {
	if level <= LvTrace {
		root.Output(2, "[TRACE] "+fmt.Sprintf(f, a...))
	}
}

// Debug ...
func Debug(a ...interface{}) {
	if level <= LvDebug {
		root.Output(2, "[DEBUG] "+fmt.Sprint(a...))
	}
}

// Debugf ...
func Debugf(f string, a ...interface{}) {
	if level <= LvDebug {
		root.Output(2, "[DEBUG] "+fmt.Sprintf(f, a...))
	}
}

// Info ...
func Info(a ...interface{}) {
	if level <= LvInfo {
		root.Output(2, "[INFO] "+fmt.Sprint(a...))
	}
}

// Infof ...
func Infof(f string, a ...interface{}) {
	if level <= LvInfo {
		root.Output(2, "[INFO] "+fmt.Sprintf(f, a...))
	}
}

// Warn ...
func Warn(a ...interface{}) {
	if level <= LvWarn {
		root.Output(2, "[WARN] "+fmt.Sprint(a...))
	}
}

// Warnf ...
func Warnf(f string, a ...interface{}) {
	if level <= LvWarn {
		root.Output(2, "[WARN] "+fmt.Sprintf(f, a...))
	}
}

// Error ...
func Error(a ...interface{}) {
	if level <= LvError {
		root.Output(2, "[ERROR] "+fmt.Sprint(a...))
	}
}

// Errorf ...
func Errorf(f string, a ...interface{}) {
	if level <= LvError {
		root.Output(2, "[Error] "+fmt.Sprintf(f, a...))
	}
}

// Fatal ...
func Fatal(a ...interface{}) {
	if level <= LvFatal {
		root.Output(2, "[FATAL] "+fmt.Sprint(a...))
	}
}

// Fatalf ...
func Fatalf(f string, a ...interface{}) {
	if level <= LvFatal {
		root.Output(2, "[FATAL] "+fmt.Sprintf(f, a...))
	}
}
