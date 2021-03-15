package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	rawLogger, _ := test.NewNullLogger()

	SetLogger(rawLogger)
	logger := GetLogger()

	assert.Equal(t, rawLogger, logger, "allows logger overriding")
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	assert.IsType(t, &logrus.Logger{}, logger, "is working without setting logger")
}

func TestLog(t *testing.T) {
	rawLogger, hook := test.NewNullLogger()
	SetLogger(rawLogger)

	Log("filereader", "formatted string")

	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level, "should have correct level")
	assert.Equal(t, "formatted string\n", hook.LastEntry().Message, "should log properly")
	assert.Equal(t, "filereader", hook.LastEntry().Data["module"], "should have correct fields set")
	hook.Reset()
}

func TestDebug(t *testing.T) {
	rawLogger, hook := test.NewNullLogger()

	rawLogger.Level = logrus.DebugLevel

	SetLogger(rawLogger)
	Debug("filereader", "formatted string")

	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level, "should have correct level")
	assert.Equal(t, "formatted string\n", hook.LastEntry().Message, "should log properly")
	assert.Equal(t, "filereader", hook.LastEntry().Data["module"], "should have correct fields set")
	hook.Reset()
}

func TestWarn(t *testing.T) {
	rawLogger, hook := test.NewNullLogger()

	rawLogger.Level = logrus.WarnLevel

	SetLogger(rawLogger)
	Warn("filereader", "formatted string")

	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level, "should have correct level")
	assert.Equal(t, "formatted string\n", hook.LastEntry().Message, "should log properly")
	assert.Equal(t, "filereader", hook.LastEntry().Data["module"], "should have correct fields set")
	hook.Reset()
}

func TestError(t *testing.T) {
	rawLogger, hook := test.NewNullLogger()

	rawLogger.Level = logrus.ErrorLevel

	SetLogger(rawLogger)
	Error("filereader", "formatted string")

	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level, "should have correct level")
	assert.Equal(t, "formatted string\n", hook.LastEntry().Message, "should log properly")
	assert.Equal(t, "filereader", hook.LastEntry().Data["module"], "should have correct fields set")
	hook.Reset()
}
