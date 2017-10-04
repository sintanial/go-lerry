package lerry

import (
	"reflect"
	"github.com/ansel1/merry"
	"github.com/mgutz/logxi/v1"
)

const NestedError = "nested error"
const UserMessage = "user message"

func Wrap(e error, level int) merry.Error {
	if e == nil {
		return nil
	}

	return merry.WrapSkipping(e, 1).WithValue("level", level)
}

func AlertWrap(e error) merry.Error {
	return Wrap(e, log.LevelAlert)
}

func FatalWrap(e error) merry.Error {
	return Wrap(e, log.LevelFatal)
}
func ErrorWrap(e error) merry.Error {
	return Wrap(e, log.LevelError)
}
func WarnWrap(e error) merry.Error {
	return Wrap(e, log.LevelWarn)
}

func NoticeWrap(e error) merry.Error {
	return Wrap(e, log.LevelNotice)
}

func InfoWrap(e error) merry.Error {
	return Wrap(e, log.LevelInfo)
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
