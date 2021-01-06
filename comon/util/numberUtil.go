package util

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

// 判断一个字符串是否为整数
func IsInteger(str string) bool {
	if len(str) == 0 {
		return false
	}
	if str[0] == '0' && len(str) > 1 {
		return false
	}
	var hasSign bool = str[0] == '-' || str[0] == '+'
	if hasSign && ((len(str) > 2 && str[1] == '0') || len(str) == 1) {
		return false
	}
	for index, ch := range str {
		if !unicode.IsDigit(ch) && !(hasSign && index == 0) {
			return false
		}
	}
	return true
}

/**
 * 保留 n 位小数，四舍五入
 * precision: 小数位
 */
func RoundFloat32(num float32, precision int) float32 {
	// 将之转为字符串，再转为指定位数小数
	f, err := strconv.ParseFloat(strconv.FormatFloat(float64(num), 'f', precision, 32), 32)
	if err != nil {
		fmt.Println(err)
	}
	return float32(f)
}

/**
 * 保留 n 位小数，四舍五入
 * precision: 小数位
 */
func RoundFloat64(num float64, precision int) float64 {
	// 将之转为字符串，再转为指定位数小数
	f, err := strconv.ParseFloat(strconv.FormatFloat(num, 'f', precision, 64), 64)
	if err != nil {
		fmt.Println("internal err in RoundFloat64:", err)
	}
	return f
}

/**
 * 将 byte 最多转为 MB
 */
const mbs = 1024 * 1024

func ByteToMB(b int) string {
	s := ""
	if b < 1024 {
		return strconv.Itoa(b) + "B"
	} else if b < mbs {
		if b%1024 > 0 {
			s = ByteToMB(b % 1024)
		}
		return strconv.Itoa(b/1024) + "KB" + s
	} else {
		if b%mbs > 0 {
			s = ByteToMB(b % mbs)
		}
		return strconv.Itoa(b/mbs) + "MB" + s
	}
}

func BytesToBinaryString(bs []byte) string {
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}

func BytesToByteString(bs []byte) []string {
	byteArr := []string{}
	for _, v := range bs {
		byteArr = append(byteArr, fmt.Sprintf("%d", v))
	}
	return byteArr
}
