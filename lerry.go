package lerry

import (
	"github.com/ansel1/merry"
	"github.com/sirupsen/logrus"
)

const KeyLogLevel = "loglevel"
const KeyLogEntry = "logentry"

func WithLevel(e error, level logrus.Level) merry.Error { return merry.WithValue(e, KeyLogLevel, level) }
func NewLevel(m string, level logrus.Level) merry.Error { return WithLevel(merry.New(m), level) }

func Level(e error) logrus.Level {
	val := merry.Value(e, KeyLogLevel)
	if val == nil {
		return 0
	}

	lvl, ok := val.(logrus.Level)
	if !ok {
		return 0
	}

	return lvl
}

func WithEntry(e error, t *logrus.Entry) merry.Error { return merry.WithValue(e, KeyLogEntry, t) }
func NewEntry(m string, t *logrus.Entry) merry.Error { return WithEntry(merry.New(m), t) }

func Entry(e error) *logrus.Entry {
	val := merry.Value(e, KeyLogEntry)
	if val == nil {
		return nil
	}

	entry, ok := val.(*logrus.Entry)
	if !ok {
		return nil
	}

	return entry
}

func WrapPanic(e error) merry.Error                       { return WithLevel(e, logrus.PanicLevel) }
func WrapPanicEntry(e error, t *logrus.Entry) merry.Error { return WithEntry(WrapPanic(e), t) }
func NewPanic(m string) merry.Error                       { return WrapPanic(merry.New(m)) }
func NewPanicEntry(m string, e *logrus.Entry) merry.Error { return WrapPanicEntry(merry.New(m), e) }

func WrapFatal(e error) merry.Error                       { return WithLevel(e, logrus.FatalLevel) }
func WrapFatalEntry(e error, t *logrus.Entry) merry.Error { return WithEntry(WrapFatal(e), t) }
func NewFatal(m string) merry.Error                       { return WrapFatal(merry.New(m)) }
func NewFatalEntry(m string, e *logrus.Entry) merry.Error { return WrapFatalEntry(merry.New(m), e) }

func WrapError(e error) merry.Error                       { return WithLevel(e, logrus.ErrorLevel) }
func WrapErrorEntry(e error, t *logrus.Entry) merry.Error { return WithEntry(WrapError(e), t) }
func NewError(m string) merry.Error                       { return WrapError(merry.New(m)) }
func NewErrorEntry(m string, e *logrus.Entry) merry.Error { return WrapErrorEntry(merry.New(m), e) }

func WrapWarn(e error) merry.Error                       { return WithLevel(e, logrus.WarnLevel) }
func WrapWarnEntry(e error, t *logrus.Entry) merry.Error { return WithEntry(WrapWarn(e), t) }
func NewWarn(m string) merry.Error                       { return WrapWarn(merry.New(m)) }
func NewWarnEntry(m string, e *logrus.Entry) merry.Error { return WrapWarnEntry(merry.New(m), e) }

func WrapDebug(e error) merry.Error                       { return WithLevel(e, logrus.DebugLevel) }
func WrapDebugEntry(e error, t *logrus.Entry) merry.Error { return WithEntry(WrapDebug(e), t) }
func NewDebug(m string) merry.Error                       { return WrapDebug(merry.New(m)) }
func NewDebugEntry(m string, e *logrus.Entry) merry.Error { return WrapDebugEntry(merry.New(m), e) }

func WrapTrace(e error) merry.Error                       { return WithLevel(e, logrus.TraceLevel) }
func WrapTraceEntry(e error, t *logrus.Entry) merry.Error { return WithEntry(WrapTrace(e), t) }
func NewTrace(m string) merry.Error                       { return WrapTrace(merry.New(m)) }
func NewTraceEntry(m string, e *logrus.Entry) merry.Error { return WrapTraceEntry(merry.New(m), e) }

func Log(err error, args ...interface{}) {
	LogWithEntry(err, Entry(err), args)
}

func LogWithEntry(err error, entry *logrus.Entry, args ...interface{}) {
	if err == nil || entry == nil {
		return
	}

	lvl := Level(err)
	if lvl == 0 {
		lvl = logrus.WarnLevel
	}

	msg := merry.UserMessage(err)
	if msg != "" {
		args = append([]interface{}{msg}, args...)
	}

	stack := merry.Stacktrace(err)
	if stack != "" {
		args = append(args, "stack", stack)
	}

	entry.WithError(err).Log(lvl, args...)
}
