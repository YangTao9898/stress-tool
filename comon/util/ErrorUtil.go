package util

import (
	"errors"
	"fmt"
)

func NewErrorf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a))
}

func WrapError(wrapMsg string, error error) error {
	return errors.New(wrapMsg + ":\n\t" + error.Error())
}
