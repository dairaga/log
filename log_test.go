package log_test

import (
	"os"
	"testing"

	"github.com/dairaga/log"
)

func mytest() {
	log.Info("A")
}

func TestMain(m *testing.M) {

	log.Trace("test")
	mytest()
	os.Exit(m.Run())
}
