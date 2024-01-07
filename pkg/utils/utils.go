// Package utils is used for functions that are used in multiple locations
package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

// LogErrorDefaultFormat is used to print errors in a uniform way
func LogErrorDefaultFormat(packageName, callerName string, err error, action string) {
	Logger.WithFields(log.Fields{"package": packageName, "method": callerName, "error": err}).Error(action)
}

// LogPanicDefaultFormat is used to print errors in a uniform way, and exit with a Panic
func LogPanicDefaultFormat(packageName, callerName string, err error, action string) {
	Logger.WithFields(log.Fields{"package": packageName, "method": callerName, "error": err}).Panic(action)
}

// LogFatalDefaultFormat is used to print errors in a uniform way, and exit with a Panic
func LogFatalDefaultFormat(packageName, callerName string, err error, action string) {
	Logger.WithFields(log.Fields{"package": packageName, "method": callerName, "error": err}).Fatal(action)
}

// NewLogger follows this stack overflow post:
// https://stackoverflow.com/a/52923899/15410622
func NewLogger(lvl log.Level) *log.Logger {
	logger := &log.Logger{
		Out:   os.Stdout,
		Level: lvl,
		Formatter: &log.TextFormatter{
			TimestampFormat: "2009-10-31T01:48:52Z",
			PadLevelText:    true,
		},
	}
	Logger = logger
	return Logger
}
