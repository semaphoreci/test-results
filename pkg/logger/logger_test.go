package logger

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func Test_SetLogger(t *testing.T) {
	logger, _ := test.NewNullLogger()
	LogEntry.SetLogger(logger)

	assert.Equal(t, logger, LogEntry.logger, "Should properly set logger")
}

func Test_GetLogger(t *testing.T) {
	logger, _ := test.NewNullLogger()
	LogEntry.SetLogger(logger)

	assert.Equal(t, logger, LogEntry.GetLogger(), "Should properly get logger")
}

func Test_Debug(t *testing.T) {
	logger, hook := test.NewNullLogger()
	LogEntry.SetLogger(logger)
	LogEntry.SetLevel(DebugLevel)

	Debug(Fields{"foo": "bar", "bar": "foo"}, "Debug")
	assert.Equal(t, "Debug\n", hook.LastEntry().Message)
	assert.Equal(t, Fields{"foo": "bar", "bar": "foo"}, Fields(hook.LastEntry().Data))
}

func Test_Warn(t *testing.T) {
	logger, hook := test.NewNullLogger()
	LogEntry.SetLogger(logger)
	LogEntry.SetLevel(WarnLevel)

	Warn(Fields{"foo": "bar", "bar": "foo"}, "Warn")
	assert.Equal(t, "Warn\n", hook.LastEntry().Message)
	assert.Equal(t, Fields{"foo": "bar", "bar": "foo"}, Fields(hook.LastEntry().Data))
}

func Test_Error(t *testing.T) {
	logger, hook := test.NewNullLogger()
	LogEntry.SetLogger(logger)
	LogEntry.SetLevel(ErrorLevel)

	Error(Fields{"foo": "bar", "bar": "foo"}, "Error")
	assert.Equal(t, "Error\n", hook.LastEntry().Message)
	assert.Equal(t, Fields{"foo": "bar", "bar": "foo"}, Fields(hook.LastEntry().Data))
}

func Test_Info(t *testing.T) {
	logger, hook := test.NewNullLogger()
	LogEntry.SetLogger(logger)
	LogEntry.SetLevel(InfoLevel)

	Info(Fields{"foo": "bar", "bar": "foo"}, "Info")
	assert.Equal(t, "Info\n", hook.LastEntry().Message)
	assert.Equal(t, Fields{"foo": "bar", "bar": "foo"}, Fields(hook.LastEntry().Data))
}

func Test_Log(t *testing.T) {
	levels := []Level{InfoLevel, WarnLevel, ErrorLevel, DebugLevel, TraceLevel}

	testCases := []struct {
		desc   string
		msg    string
		fields Fields
	}{
		{
			desc:   "Works in info Level",
			msg:    "Testing logs ...",
			fields: Fields{"foo": "bar", "bar": "foo"},
		},
	}

	logger, hook := test.NewNullLogger()
	LogEntry.SetLogger(logger)

	for _, tC := range testCases {
		for _, level := range levels {
			t.Run(tC.desc+" on level "+fmt.Sprint(level), func(t *testing.T) {
				LogEntry.SetLevel(level)
				LogEntry.Log(level, tC.fields, tC.msg)
				assert.Equal(t, tC.msg+"\n", hook.LastEntry().Message)
				assert.Equal(t, tC.fields, Fields(hook.LastEntry().Data))

				hook.Reset()
			})
		}
	}
}
