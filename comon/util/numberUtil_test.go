package util

import (
	"strconv"
	"testing"
)

func TestIsInteger(t *testing.T) {
	testStrMap := map[string]bool{ // 测试用例，key 为测试字符串，value 为结果
		"0":    true,
		"01":   false,
		"10":   true,
		"-":    false,
		"-10":  true,
		"-01":  false,
		"+":    false,
		"+10":  true,
		"+01":  false,
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
	mb = ByteToMB(1024*1024 + 1)
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

func TestBytesToBinaryString(t *testing.T) {
	type testcase struct {
		bytes     []byte
		expectRes string
	}

	testcases := []testcase{}

	testcases = append(testcases, testcase{
		bytes:     []byte{1},
		expectRes: "00000001",
	})

	testcases = append(testcases, testcase{
		bytes:     []byte{0, 1},
		expectRes: "0000000000000001",
	})

	testcases = append(testcases, testcase{
		bytes:     []byte{255, 254},
		expectRes: "1111111111111110",
	})

	testcases = append(testcases, testcase{
		bytes:     []byte{0, 2, 0},
		expectRes: "000000000000001000000000",
	})

	for _, v := range testcases {
		res := BytesToBinaryString(v.bytes, false)
		t.Logf("testcase [%+v] expect [%s], get [%s]", v.bytes, v.expectRes, res)
		if res != v.expectRes {
			t.Error("fail")
		}
	}
}

func TestBinaryStringToBytes(t *testing.T) {
	inbs := []byte{123, 127, 127}
	bstr := BytesToBinaryString(inbs, true)
	t.Logf("bstr: %s", bstr)
	bs, err := BinaryStringToBytes(bstr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", bs)
}

func TestGetUniqueString(t *testing.T) {
	for i := 0; i < 30; i++ {
		s, err := GetUniqueString("test-")
		if err != nil {
			t.Errorf("err:[%+v], result:[%s]", err, s)
		} else {
			t.Log(s)
		}
	}
}
