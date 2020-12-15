package util

import "unicode"

// 判断一个字符串是否为整数
func IsInteger(str string) bool {
	if len(str) == 0 {
		return false
	}
	if str[0] == '0' && len(str) > 1 {
		return false
	}
	var hasSign bool = str[0] == '-' || str[0] == '+'
	if  hasSign && ((len(str) > 2 && str[1] == '0') || len(str) == 1) {
		return false
	}
	for index, ch := range str {
		if !unicode.IsDigit(ch) && !(hasSign && index == 0) {
			return false;
		}
	}
	return true
}
