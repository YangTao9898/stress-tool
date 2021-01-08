package tcp

import (
	"golang.org/x/tools/container/intsets"
	"sort"
	"strconv"
	"stress-tool/model"
)

type countData struct {
	totalRequestCount int
	totalCostTime     float64
	costMaxTime       float64
	costMinTime       float64
	//totalResponseTime float64
	//responseMaxTime float64
	//responseMinTime float64
	threadCostTime    int64
	threadCostMaxTime int64
	threadCostMinTIme int64
	succCount         int
	failCount         int
	timeoutCount      int
	sendDataSize      int // 发送的数据量，单位 byte
	reciveDataSize    int
	isFail            bool // 此次连接是否失败
	isTimeout         bool // 此次连接是否超时
	isWaitReadEnd     bool // 未填写期望返回字节时，是否有正常等待结束的时候
}

// 检查数据范围是否连续
// return: int 为范围最大值 string 为返回的错误信息
func checkDataRangeIsContinuous(mArr []map[string]int) (int, string) {
	rangelistMap := make(map[int][2]int)
	var L = make([]int, len(mArr))
	index := 0

	m := make(map[string]int, len(mArr))
	// 合并多个Map
	for _, v := range mArr {
		for k, o := range v {
			m[k] = o
		}
	}
	for key, value := range m {
		n1, n2, err := model.CreateTaskDataKeySplite(key)
		if err != nil {
			return -1, err.Error()
		}
		arr := [2]int{n1, n2}
		if arr[0] > arr[1] {
			return -1, "数据范围开始索引不能大于结束索引"
		}
		if value < 0 || value > 2 {
			return -1, "未知的数据类型"
		}
		dataLength := arr[1] - arr[0] + 1
		if value == 0 && !dataIntLength[dataLength] {
			return -1, "INT类型长度不能为[" + strconv.Itoa(dataLength) + "], 必须在1,2,4,8之间"
		}
		if value == 1 && !dataFloatLength[dataLength] {
			return -1, "FLOAT类型长度不能为[" + strconv.Itoa(dataLength) + "]，必须在4,8之间"
		}

		rangelistMap[arr[0]] = arr
		L[index] = arr[0]
		index++
	}
	// 测试 key 对应的 value 所表示的范围是否连续
	sort.Ints(L)
	if L[0] != 0 { // 下标不是从0开始的
		return -1, "数据范围不是从0开始"
	}
	var arr [2]int
	for index, value := range L {
		tempArr := rangelistMap[value]
		if index == 0 {
			arr = tempArr
			continue
		}
		if arr[1]+1 != tempArr[0] {
			return -1, "数据范围不连续"
		}
		arr = tempArr
	}
	// 取出索引范围内最大的值, 加 1 得到数据长度
	return rangelistMap[L[len(L)-1]][1] + 1, ""
}

func checkCreateTaskData(data model.CreateTaskData) string {
	if data.ThreadNum <= 0 {
		return "线程数不能小于等于0"
	}
	if data.DataTypeMap == nil {
		return "请求数据类型不能为空"
	}
	if data.HasResponse && data.ReadTimeout <= 0 {
		return "读取超时时间必须大于0"
	}

	length, msg := checkDataRangeIsContinuous(data.DataTypeMap)
	if msg != "" {
		return msg
	}
	if len(data.Data) > length {
		return "定义的数据长度小于实际的数据长度"
	}
	if data.IsRepeat {
		if data.RepeatTime <= 0 {
			return "重复请求次数须大于0"
		}
		if data.IntervalTime < 0 {
			return "请求间隔不能小于0"
		}
	}

	return ""
}

func mergeSingleResult(cdata countData, totalCount *countData, isSucc bool) {
	if isSucc {
		totalCount.succCount += cdata.succCount
		if totalCount.costMaxTime < cdata.totalCostTime {
			totalCount.costMaxTime = cdata.totalCostTime
		}
		if totalCount.costMinTime > cdata.totalCostTime {
			totalCount.costMinTime = cdata.totalCostTime
		}
		/*if totalCount.responseMaxTime < cdata.totalResponseTime {
			totalCount.responseMaxTime = cdata.totalResponseTime
		}
		if totalCount.responseMinTime > cdata.totalResponseTime {
			totalCount.responseMinTime = cdata.totalResponseTime
		}*/
		totalCount.sendDataSize += cdata.sendDataSize
		totalCount.reciveDataSize += cdata.reciveDataSize
		if cdata.isWaitReadEnd {
			totalCount.isWaitReadEnd = true
		}
	}
	totalCount.timeoutCount += cdata.timeoutCount
	totalCount.failCount += cdata.failCount
	totalCount.totalRequestCount = totalCount.succCount + totalCount.timeoutCount + totalCount.failCount
}

func mergeThreadResult(cdata countData, totalCount *countData, isSucc bool) {
	if isSucc {
		totalCount.succCount += cdata.succCount
		//totalCount.totalResponseTime += cdata.totalResponseTime
		//totalCount.totalCostTime += cdata.totalCostTime
		if totalCount.costMaxTime < cdata.costMaxTime {
			totalCount.costMaxTime = cdata.costMaxTime
		}
		if totalCount.costMinTime > cdata.costMinTime {
			totalCount.costMinTime = cdata.costMinTime
		}
		/*if totalCount.responseMaxTime < cdata.responseMaxTime {
			totalCount.responseMaxTime = cdata.responseMaxTime
		}
		if totalCount.responseMinTime > cdata.responseMinTime {
			totalCount.responseMinTime = cdata.responseMinTime
		}*/
		totalCount.sendDataSize += cdata.sendDataSize
		totalCount.reciveDataSize += cdata.reciveDataSize
		if cdata.isWaitReadEnd {
			totalCount.isWaitReadEnd = true
		}
	}
	totalCount.threadCostTime += cdata.threadCostTime
	if totalCount.threadCostMaxTime < cdata.threadCostTime {
		totalCount.threadCostMaxTime = cdata.threadCostTime
	}
	if totalCount.threadCostMinTIme > cdata.threadCostTime {
		totalCount.threadCostMinTIme = cdata.threadCostTime
	}

	totalCount.timeoutCount += cdata.timeoutCount
	totalCount.failCount += cdata.failCount
	totalCount.totalRequestCount = totalCount.succCount + totalCount.timeoutCount + totalCount.failCount
}

func getDefaultTotalCData() countData {
	return countData{
		totalRequestCount: 0,
		totalCostTime:     0,
		costMaxTime:       float64(intsets.MinInt),
		costMinTime:       float64(intsets.MaxInt),
		/*totalResponseTime: 0,
		responseMaxTime:   float64(intsets.MinInt),
		responseMinTime:   float64(intsets.MaxInt),*/
		threadCostTime:    0,
		threadCostMaxTime: int64(intsets.MinInt),
		threadCostMinTIme: int64(intsets.MaxInt),
		succCount:         0,
		failCount:         0,
		timeoutCount:      0,
		sendDataSize:      0,
		reciveDataSize:    0,
	}
}
