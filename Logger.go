package log

import (
	"encoding/json"
	"io"
	"time"
)

// Logger ...
type Logger interface {
	Output(now time.Time, severity Severity, file string, lineNo int, message string)
	OutputJSON(now time.Time, severity Severity, file string, lineNo int, message json.RawMessage)
	OutputStruct(now time.Time, severity Severity, file string, lineNo int, message interface{})
	io.Closer
}
