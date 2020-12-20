package util

import (
	"strconv"
	"testing"
)

func TestIsInteger(t *testing.T) {
	testStrMap := map[string]bool{ // 测试用例，key 为测试字符串，value 为结果
		"0": true,
		"01": false,
		"10": true,
		"-": false,
		"-10": true,
		"-01": false,
		"+": false,
		"+10": true,
		"+01": false,
		"100a": false,
		"10a0": false,
		"a100": false,
	}
	for key, value := range testStrMap {
		b := IsInteger(key)
		t.Logf("test IsInteger, param: [%s] expected: [%t], get [%t]", key, value, b)
		if b != value {
			t.Error("IsInteger failed")
		}
	}
}

func TestRoundFloat32(t *testing.T) {
	var f float32 = 3.123456789
	res := RoundFloat32(f, 2)
	t.Logf("expect: [3.12] get: [%s]", strconv.FormatFloat(float64(res), 'f', -1, 32))
	if res != 3.12 {
		t.Error("not expect")
	}

	f = 3.125456789
	res = RoundFloat32(f, 2)
	t.Logf("expect: [3.13] get: [%s]", strconv.FormatFloat(float64(res), 'f', -1, 32))
	if res != 3.13 {
		t.Error("not expect")
	}
}

func TestRoundFloat64(t *testing.T) {
	var f float64 = 3.123456789
	res := RoundFloat64(f, 2)
	t.Logf("expect: [3.12] get: [%s]", strconv.FormatFloat(res, 'f', -1, 64))
	if res != 3.12 {
		t.Error("not expect")
	}

	f = 3.125456789
	res = RoundFloat64(f, 2)
	t.Logf("expect: [3.13] get: [%s]", strconv.FormatFloat(res, 'f', -1, 64))
	if res != 3.13 {
		t.Error("not expect")
	}
}

func TestByteToMB(t *testing.T) {
	expect := "0B"
	mb := ByteToMB(0)
	t.Logf("expect: [%s], get: [%s]", expect, mb)
	if mb != expect {
		t.Error("fail")
	}

	expect = "1023B"
	mb = ByteToMB(1023)
	t.Logf("expect: [%s], get: [%s]", expect, mb)
	if mb != expect {
		t.Error("fail")
	}

	expect = "1KB"
	mb = ByteToMB(1024)
	t.Logf("expect: [%s], get: [%s]", expect, mb)
	if mb != expect {
		t.Error("fail")
	}

	expect = "1KB1B"
	mb = ByteToMB(1025)
	t.Logf("expect: [%s], get: [%s]", expect, mb)
	if mb != expect {
		t.Error("fail")
	}

	expect = "1MB"
	mb = ByteToMB(1024 * 1024)
	t.Logf("expect: [%s], get: [%s]", expect, mb)
	if mb != expect {
		t.Error("fail")
	}

	expect = "1MB1B"
	mb = ByteToMB(1024 * 1024 + 1)
	t.Logf("expect: [%s], get: [%s]", expect, mb)
	if mb != expect {
		t.Error("fail")
	}

	expect = "1024MB"
	mb = ByteToMB(1024 * 1024 * 1024)
	t.Logf("expect: [%s], get: [%s]", expect, mb)
	if mb != expect {
		t.Error("fail")
	}
}