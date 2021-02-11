package util

import (
	"errors"
	"fmt"
	"runtime"
)

func NewErrorf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

func WrapError(wrapMsg string, error error) error {
	pc, file, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return errors.New("WrapError 方法获取堆栈失败")
	}
	if error == nil {
		return nil
	} else {
		errMsg := fmt.Sprintf("%s \n\tat %s:%d (Method %s)\nCause by: %s\n", wrapMsg, file, line, f.Name(), error.Error())
		return errors.New(errMsg)
	}
}
