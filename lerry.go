package lerry

import (
	"github.com/ansel1/merry"
	"github.com/sirupsen/logrus"
)

const KeyLogLevel = "loglevel"

func wrap(e error, level logrus.Level) merry.Error {
	if e == nil {
		return nil
	}

	return merry.WrapSkipping(e, 3).WithValue(KeyLogLevel, level)
}

func level(e error) logrus.Level {
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

func PanicWrap(e error) merry.Error {
	return wrap(e, logrus.PanicLevel)
}

func FatalWrap(e error) merry.Error {
	return wrap(e, logrus.FatalLevel)
}

func ErrorWrap(e error) merry.Error {
	return wrap(e, logrus.ErrorLevel)
}

func WarnWrap(e error) merry.Error {
	return wrap(e, logrus.WarnLevel)
}

func InfoWrap(e error) merry.Error {
	return wrap(e, logrus.InfoLevel)
}

func DebugWrap(e error) merry.Error {
	return wrap(e, logrus.DebugLevel)
}

func TraceWrap(e error) merry.Error {
	return wrap(e, logrus.TraceLevel)
}

func Log(err error, entry *logrus.Entry, args ...interface{}) {
	lvl := level(err)
	if lvl == 0 {
		lvl = logrus.DebugLevel
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
