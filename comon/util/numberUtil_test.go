package util

import "testing"

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
