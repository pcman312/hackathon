package main

import (
	log "github.com/cihub/seelog"
)

// ESInfoLogger logs to seelog as an info message
type ESInfoLogger struct{}

// Printf ...
func (ESInfoLogger) Printf(format string, v ...interface{}) {
	log.Infof(format, v...)
}

// ESTraceLogger logs to seelog as a trace message
type ESTraceLogger struct{}

// Printf ...
func (ESTraceLogger) Printf(format string, v ...interface{}) {
	log.Tracef(format, v...)
}

// ESErrorLogger logs to seelog as an error message
type ESErrorLogger struct{}

// Printf ...
func (ESErrorLogger) Printf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}
