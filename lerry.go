package lerry

import (
	"reflect"
	"github.com/ansel1/merry"
	"github.com/mgutz/logxi/v1"
)

const NestedError = "nested error"
const UserMessage = "user message"

func wrap(e error, level int, skipstack int) merry.Error {
	if e == nil {
		return nil
	}

	return merry.WrapSkipping(e, skipstack).WithValue("level", level)
}

func AlertWrap(e error) merry.Error {
	return wrap(e, log.LevelAlert, 2)
}

func FatalWrap(e error) merry.Error {
	return wrap(e, log.LevelFatal, 2)
}
func ErrorWrap(e error) merry.Error {
	return wrap(e, log.LevelError, 2)
}
func WarnWrap(e error) merry.Error {
	return wrap(e, log.LevelWarn, 2)
}

func NoticeWrap(e error) merry.Error {
	return wrap(e, log.LevelNotice, 2)
}

func InfoWrap(e error) merry.Error {
	return wrap(e, log.LevelInfo, 2)
}

func Print(lg log.Logger, err error) {
	level, msg, args := Prepare(err)
	lg.Log(level, msg, args)
}

func NestedPrint(lg log.Logger, err error) {
	if nerr, ok := merry.Value(err, NestedError).(error); ok {
		NestedPrint(lg, nerr)
	}
	Print(lg, err)
}

func Prepare(err error) (level int, message string, arguments []interface{}) {
	values := merry.Values(err)

	args := make([]interface{}, 0)

	for key, val := range values {
		if key == UserMessage || key == NestedError {
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
		lv = log.LevelWarn
	}

	return lv, usrmsg + ": " + errmsg, args
}
