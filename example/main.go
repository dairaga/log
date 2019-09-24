package main

import (
	"time"

	"github.com/dairaga/log"
)

func mytest() {
	log.Info("A")
}

type mydata struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func main() {
	log.Debug("debug 哈哈哈")
	log.Debugf("debug aa: %s", "test")

	log.Fatal("xx c c c")
	log.Fatalf("xxx cc: %s", "test")

	mytest()
	log.Close(1 * time.Second)
}
