package util

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrapError(t *testing.T) {
	err := errors.New("you are wrong")
	err = WrapError("oneWrap", err)
	err = WrapError("twoWrap", err)
	fmt.Print(err)
}

func TestGetErrorStack(t *testing.T) {
	err := errors.New("I am an error")
	fmt.Println(GetErrorStack("real an error", err))
}
