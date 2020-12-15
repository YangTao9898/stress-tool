package util

import (
	"time"
)

const (
	TIME_TEMPLATE_1 string = "2006-01-02 15:04:05"
	TIME_TEMPLATE_2 string = "20060102150405"
)

func GetNowDateToStr(timeTemplate string) string{
	return time.Now().Format(timeTemplate)
}

func GetDateToStr(t time.Time, timeTemplate string) string {
	return t.Format(timeTemplate)
}

// 生成时间戳，单位ms
func GetNowTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// 生成时间戳，单位ms
func GetTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}





