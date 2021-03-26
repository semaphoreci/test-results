package logger

import (
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// LogEntry ...
var LogEntry *logEntry = new()

// Level ...
type Level uint32

// Fields ...
type Fields map[string]interface{}

const (
	//InfoLevel ...
	InfoLevel = Level(log.InfoLevel)

	// ErrorLevel ...
	ErrorLevel = Level(log.ErrorLevel)

	// WarnLevel ...
	WarnLevel = Level(log.WarnLevel)

	// DebugLevel ...
	DebugLevel = Level(log.DebugLevel)

	// TraceLevel ...
	TraceLevel = Level(log.TraceLevel)
)

// LogEntry  ...
type logEntry struct {
	logger *log.Logger
}

// new ...
func new() *logEntry {
	return &logEntry{logger: log.New()}
}

// GetLogger ...
func (me *logEntry) GetLogger() *log.Logger {
	return me.logger
}

// SetLogger ...
func (me *logEntry) SetLogger(logger *log.Logger) {
	(*me).logger = logger
}

// SetLevel ...
func (me *logEntry) SetLevel(level Level) {
	(*me).logger.SetLevel(log.Level(level))
}

// Error ...
func Error(fields Fields, s string, args ...interface{}) {
	LogEntry.Log(ErrorLevel, fields, s, args...)
}

// Warn ...
func Warn(fields Fields, s string, args ...interface{}) {
	LogEntry.Log(WarnLevel, fields, s, args...)
}

// Info ...
func Info(fields Fields, s string, args ...interface{}) {
	LogEntry.Log(InfoLevel, fields, s, args...)
}

// Debug ...
func Debug(fields Fields, s string, args ...interface{}) {
	LogEntry.Log(DebugLevel, fields, s, args...)
}

// Trace ...
func Trace(fields Fields, s string, args ...interface{}) {
	LogEntry.Log(TraceLevel, fields, s, args...)
}

// Log ...
// TODO: change to separate methods (?)
func (me *logEntry) Log(level Level, fields Fields, s string, args ...interface{}) {
	switch level {
	case ErrorLevel:
		me.logger.WithFields(logrus.Fields(fields)).Errorf(s+"\n", args...)
	case WarnLevel:
		me.logger.WithFields(logrus.Fields(fields)).Warnf(s+"\n", args...)
	case InfoLevel:
		me.logger.WithFields(logrus.Fields(fields)).Infof(s+"\n", args...)
	case DebugLevel:
		me.logger.WithFields(logrus.Fields(fields)).Debugf(s+"\n", args...)
	case TraceLevel:
		me.logger.WithFields(logrus.Fields(fields)).Tracef(s+"\n", args...)
	}
}
