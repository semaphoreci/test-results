package logger

import (
	log "github.com/sirupsen/logrus"
)

var logger = log.New()

// SetLogger ...
func SetLogger(newLogger *log.Logger) {
	logger = newLogger
}

// GetLogger ...
func GetLogger() *log.Logger {
	return logger
}

// Log ...
func Log(module string, msg string, inter ...interface{}) {
	ctx := GetLogger()
	ctx.WithFields(log.Fields{
		"module": module,
	}).Printf(msg+"\n", inter...)
}

// Debug ...
func Debug(module string, msg string, inter ...interface{}) {
	ctx := GetLogger()
	ctx.WithFields(log.Fields{
		"module": module,
	}).Debugf(msg+"\n", inter...)
}

// Error ...
func Error(module string, msg string, inter ...interface{}) {
	ctx := GetLogger()
	ctx.WithFields(log.Fields{
		"module": module,
	}).Errorf(msg+"\n", inter...)
}

// Warn ...
func Warn(module string, msg string, inter ...interface{}) {
	ctx := GetLogger()
	ctx.WithFields(log.Fields{
		"module": module,
	}).Warnf(msg+"\n", inter...)
}
