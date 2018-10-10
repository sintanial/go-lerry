package lerry

import (
	"reflect"
	"github.com/ansel1/merry"
)

type LoggerFunc func(level int, msg string, args []interface{})

const (
	// LevelAlert means action must be taken immediately.
	LevelAlert = 1

	// LevelFatal means it should be corrected immediately, eg cannot connect to database.
	LevelFatal = 2

	// LevelError is a non-urgen failure to notify devlopers or admins
	LevelError = 3

	// LevelWarn indiates an error will occur if action is not taken, eg file system 85% full
	LevelWarn = 4

	// LevelNotice is normal but significant condition.
	LevelNotice = 5

	// LevelInfo is info level
	LevelInfo = 6

	// LevelDebug is debug level
	LevelDebug = 7

	// LevelTrace is trace level and displays file and line in terminal
	LevelTrace = 8
)

const KeyNestedError = "nested error"
const KeyUserMessage = "user message"
const KeyLevel = "level"
const KeyHttpResponseStatusCode = "http status code"

func wrap(e error, level int) merry.Error {
	if e == nil {
		return nil
	}

	return merry.WrapSkipping(e, 3).WithValue("level", level)
}

func AlertWrap(e error) merry.Error {
	return wrap(e, LevelAlert)
}

func FatalWrap(e error) merry.Error {
	return wrap(e, LevelFatal)
}
func ErrorWrap(e error) merry.Error {
	return wrap(e, LevelError)
}
func WarnWrap(e error) merry.Error {
	return wrap(e, LevelWarn)
}

func NoticeWrap(e error) merry.Error {
	return wrap(e, LevelNotice)
}

func InfoWrap(e error) merry.Error {
	return wrap(e, LevelInfo)
}

func DebugWrap(e error) merry.Error {
	return wrap(e, LevelDebug)
}

func TraceWrap(e error) merry.Error {
	return wrap(e, LevelTrace)
}

func Print(lg LoggerFunc, err error) {
	level, msg, args := Prepare(err)
	lg(level, msg, args)
}

func NestedPrint(lg LoggerFunc, err error) {
	if nerr, ok := merry.Value(err, KeyNestedError).(error); ok {
		NestedPrint(lg, nerr)
	}
	Print(lg, err)
}

func Prepare(err error) (level int, message string, arguments []interface{}) {
	values := merry.Values(err)

	args := make([]interface{}, 0)

	for key, val := range values {
		if key == KeyUserMessage || key == KeyNestedError || key == KeyLevel || key == KeyHttpResponseStatusCode {
			continue
		}

		if reflect.TypeOf(key).String() == "string" {
			args = append(args, key, val)
		}
	}

	if strace := merry.Stacktrace(err); strace != "" {
		args = append(args, "stack", strace)
	}

	errmsg := err.Error()
	usrmsg := merry.UserMessage(err)
	lv, ok := merry.Value(err, "level").(int)
	if !ok {
		lv = LevelWarn
	}

	return lv, usrmsg + ": " + errmsg, args
}
