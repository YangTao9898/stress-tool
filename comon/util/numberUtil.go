package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

const (
	zero = byte('0')
	one  = byte('1')
)

var uint8arr [8]uint8

func init() {
	uint8arr[0] = 128
	uint8arr[1] = 64
	uint8arr[2] = 32
	uint8arr[3] = 16
	uint8arr[4] = 8
	uint8arr[5] = 4
	uint8arr[6] = 2
	uint8arr[7] = 1
}

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

/**
 * isSpaceSplite 是否每八位用空格隔开
 */
func BytesToBinaryString(bs []byte, isSpaceSplite bool) string {
	buf := bytes.NewBuffer([]byte{})
	length := len(bs)
	for i, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
		if isSpaceSplite && i < length-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

var zeroOneReg = regexp.MustCompile(`[^01]`)
var spaceReg = regexp.MustCompile("\\s+")

func BinaryStringToBytes(binaryStr string) (res []byte, retErr error) {
	binaryStr = spaceReg.ReplaceAllString(binaryStr, "")
	if len(binaryStr) == 0 {
		return res, nil
	}

	invalidStr := zeroOneReg.FindString(binaryStr)
	if "" != invalidStr {
		return res, fmt.Errorf("contains invalid char '%s'", invalidStr)
	}

	l := len(binaryStr)

	mo := l % 8
	l /= 8
	if mo != 0 {
		l++
	}
	res = make([]byte, 0, l)
	mo = 8 - mo
	var n uint8
	for i, b := range []byte(binaryStr) {
		m := (i + mo) % 8
		switch b {
		case one:
			n += uint8arr[m]
		}
		if m == 7 {
			res = append(res, n)
			n = 0
		}
	}
	return
}

func BytesToByteString(bs []byte) []string {
	byteArr := []string{}
	for _, v := range bs {
		byteArr = append(byteArr, fmt.Sprintf("%d", v))
	}
	return byteArr
}

// 生成 前缀 + 毫秒级时间戳 + 10位随机数 随机字符串
const forUniqueNum = 10000000000

func GetUniqueString(prefix string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(forUniqueNum))
	return fmt.Sprintf("%s%s-%010d", prefix, GetDateToStrWithMillisecond(time.Now()), n), err
}
