package log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"github.com/dairaga/config"
)

var (
	cliClosed = false
	cli       *logging.Client
)

type gcplogger struct {
	*logging.Logger
}

type entry struct {
	File    string      `json:"file"`
	Line    int         `json:"line"`
	Message interface{} `json:"message"`
}

// ----------------------------------------------------------------------------

func initGCPLogger() {
	c, err := config.Load("gcp_log.toml")
	if err != nil {
		return
	}
	project := c.GetString("log.project")
	name := c.GetString("log.name")
	cli, err = logging.NewClient(context.Background(), project)
	if err != nil {
		fmt.Printf("connect to gcp stackdriver logging: %v\n", err)
		return
	}

	Register(&gcplogger{cli.Logger(name)})

}

// ----------------------------------------------------------------------------------------------------------------

func (l *gcplogger) Output(now time.Time, severity Severity, file string, lineNo int, message string) {
	l.Log(
		logging.Entry{
			Timestamp: now,
			Severity:  logging.Severity(severity),
			Payload: entry{
				File:    file,
				Line:    lineNo,
				Message: message,
			},
		},
	)
}

func (l *gcplogger) OutputJSON(now time.Time, severity Severity, file string, lineNo int, message json.RawMessage) {
	l.Log(
		logging.Entry{
			Timestamp: now,
			Severity:  logging.Severity(severity),
			Payload: entry{
				File:    file,
				Line:    lineNo,
				Message: message,
			},
		},
	)
}

func (l *gcplogger) OutputStruct(now time.Time, severity Severity, file string, lineNo int, message interface{}) {
	l.Log(
		logging.Entry{
			Timestamp: now,
			Severity:  logging.Severity(severity),
			Payload: entry{
				File:    file,
				Line:    lineNo,
				Message: message,
			},
		},
	)
}

// ----------------------------------------------------------------------------------------------------------------

func (l *gcplogger) Close() error {
	if !cliClosed {
		cliClosed = true

		return cli.Close()
	}

	return nil
}
