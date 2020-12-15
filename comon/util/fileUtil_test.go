package util

import (
	"fmt"
	"testing"
)

func TestListChildWholeFileName(t *testing.T) {
	strings, e := ListChildWholeFileName("../../web-template")
	if e != nil {
		t.Error(e)
	}
	fmt.Print("[")
	for _, v := range strings {
		fmt.Print("(" + v + ")")
		fmt.Print(", ")
	}
	fmt.Print("]")
}
