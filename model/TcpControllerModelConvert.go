package model

import (
	"bytes"
	"encoding/binary"
	go_logger "github.com/phachon/go-logger"
	"strconv"
	"stress-tool/comon/util"
)

var log *go_logger.Logger

func init() {
	log = util.GetLogger()
}

/**
 * CreateTaskData 返回值
 * resCode string 错误码，为空则没错误
 */
func CheckToCreateTaskData(req CreateTaskRequest) (res CreateTaskData, resCode string) {
	if req.TargetAddress == "" {
		resCode += "1001,"
	}

	if req.TargetPort == "" {
		resCode += "1002,"
	}

	timeout, err := strconv.Atoi(req.Timeout)
	if err != nil && req.Timeout != "" { // 选填
		resCode += "1003,"
	}

	readTimeout := 0
	if req.HasResponse {
		readTimeout, err = strconv.Atoi(req.ReadTimeout)
		if err != nil {
			if req.ReadTimeout == "" {
				resCode += "1004,"
			} else {
				resCode += "1005,"
			}
		}
	}

	expectedBytes := 0
	if req.HasResponse {
		expectedBytes, err = strconv.Atoi(req.ExpectedBytes)
		if err != nil && req.ExpectedBytes != "" { // 选填
			resCode += "1007,"
		}
	}

	threadNum, err := strconv.Atoi(req.ThreadNum)
	if err != nil {
		if req.ThreadNum == "" {
			resCode += "1008,"
		} else {
			resCode += "1009,"
		}
	}

	repeatTime := 0
	if req.IsRepeat {
		repeatTime, err = strconv.Atoi(req.RepeatTime)
		if err != nil {
			if req.RepeatTime == "" {
				resCode += "1010,"
			} else {
				resCode += "1011,"
			}
		}
	}

	intervalTime := 0
	if req.IsRepeat {
		intervalTime, err = strconv.Atoi(req.IntervalTime)
		if err != nil {
			if req.IntervalTime == "" {
				resCode += "1012,"
			} else {
				resCode += "1013,"
			}
		}
	}

	if len(req.DataMapArr) == 0 {
		resCode += "1019"
		return
	}
	dataMapArr := req.DataMapArr
	databytes := bytes.NewBuffer([]byte{})
	dataTypeMap := make([]map[string]int, len(dataMapArr))
	num := 0
	for index, v := range dataMapArr {
		lengthStr := v.Length
		data := v.Data
		isBigEnd := v.IsBigEnd
		tn, err := strconv.Atoi(v.Type)
		if err != nil {
			resCode += "1080"
			return
		}

		var length int
		if tn == 2 {
			length = len(data)
		} else {
			length, err = strconv.Atoi(lengthStr)
			if err != nil {
				resCode += "1020"
				return
			}
		}

		end := num + length - 1
		m := make(map[string]int, 1)
		m[strconv.Itoa(num)+"~"+strconv.Itoa(end)] = tn
		dataTypeMap[index] = m
		num = end + 1
		switch tn {
		case NUMBER:
			switch length {
			case 1:
				dataInt, err := strconv.ParseInt(data, 10, 8)
				if err != nil {
					resCode += "1030"
					return
				}
				if isBigEnd {
					err = binary.Write(databytes, binary.BigEndian, int8(dataInt))
				} else {
					err = binary.Write(databytes, binary.BigEndian, int8(dataInt))
				}
				if err != nil {
					log.Error(err.Error())
				}
			case 2:
				dataInt, err := strconv.ParseInt(data, 10, 16)
				if err != nil {
					resCode += "1031"
					return
				}
				if isBigEnd {
					err = binary.Write(databytes, binary.BigEndian, int16(dataInt))
				} else {
					err = binary.Write(databytes, binary.BigEndian, int16(dataInt))
				}
				if err != nil {
					log.Error(err.Error())
				}
			case 4:
				dataInt, err := strconv.ParseInt(data, 10, 32)
				if err != nil {
					resCode += "1032"
					return
				}
				if isBigEnd {
					err = binary.Write(databytes, binary.BigEndian, int32(dataInt))
				} else {
					err = binary.Write(databytes, binary.BigEndian, int32(dataInt))
				}
				if err != nil {
					log.Error(err.Error())
				}
			case 8:
				dataInt, err := strconv.ParseInt(data, 10, 64)
				if err != nil {
					resCode += "1033"
					return
				}
				if isBigEnd {
					err = binary.Write(databytes, binary.BigEndian, dataInt)
				} else {
					err = binary.Write(databytes, binary.BigEndian, dataInt)
				}
				if err != nil {
					log.Error(err.Error())
				}
			default:
				resCode += "1040"
				return
			}

		case FLOAT:
			switch length {
			case 4:
				dataFloat, err := strconv.ParseFloat(data, 32)
				if err != nil {
					resCode += "1050"
					return
				}
				if isBigEnd {
					err = binary.Write(databytes, binary.BigEndian, float32(dataFloat))
				} else {
					err = binary.Write(databytes, binary.BigEndian, float32(dataFloat))
				}
				if err != nil {
					log.Error(err.Error())
				}
			case 8:
				dataFloat, err := strconv.ParseFloat(data, 64)
				if err != nil {
					resCode += "1051"
					return
				}
				if isBigEnd {
					err = binary.Write(databytes, binary.BigEndian, dataFloat)
				} else {
					err = binary.Write(databytes, binary.BigEndian, dataFloat)
				}
				if err != nil {
					log.Error(err.Error())
				}
			default:
				resCode += "1060"
			}
		case STRING:
			length = len(data)
			bytes := []byte(data)
			_, err := databytes.Write(bytes)
			if err != nil {
				log.Error(err.Error())
			}
		default:
			resCode += "1080"
			return
		}
	}

	res = CreateTaskData{
		TargetAddress: req.TargetAddress,
		TargetPort:    req.TargetPort,
		Timeout:       timeout,
		ReadTimeout:   readTimeout,
		ExpectedBytes: expectedBytes,
		ThreadNum:     threadNum,
		IsRepeat:      req.IsRepeat,
		RepeatTime:    repeatTime,
		IntervalTime:  intervalTime,
		HasResponse:   req.HasResponse,
		DataTypeMap:   dataTypeMap,
		Data:          databytes.Bytes(),
	}

	return
}

func TaskDealDataToGetTaskDetailResponse(data TaskDealData) GetTaskDetailResponse {
	readTimeout := "-"
	expectedBytes := "-"
	if data.HasResponse {
		readTimeout = strconv.Itoa(data.ReadTimeout) + "ms"
		expectedBytes = strconv.Itoa(data.ExpectedBytes) + "Byte"
	}

	repeatTime := "-"
	sendInterval := "-"
	if data.IsRepeat {
		repeatTime = strconv.Itoa(data.RepeatTime)
		sendInterval = strconv.Itoa(data.IntervalTime) + "ms"
	}

	state := ""
	switch data.State {
	case NOT_START:
		state = "未执行"
	case READY:
		state = "等待执行"
	case RUNNING:
		state = "正在执行"
	case FINISH:
		state = "执行完毕"
	default:
		state = "未知状态"
	}

	startTime := "-"
	endTime := "-"
	totalRequestCount := "-"
	requestAverageCostTime := "-"
	requestCostMaxTime := "-"
	requestCostMinTime := "-"
	requestAverageResponseTime := "-"
	requestResponseMaxTime := "-"
	requestResponseMinTime := "-"
	transactionRate := "-"
	succTransactions := "-"
	failTransactions := "-"
	timeOutTransactions := "-"
	dataTransferred := "-"
	throughput := "-"
	totalCostTime := "-"
	if data.State > READY {
		startTime = data.StartTime
	} else if data.State > RUNNING {
		startTime = data.StartTime
		endTime = data.EndTime
		totalRequestCount = strconv.Itoa(data.TotalRequestCount)
		requestAverageCostTime = strconv.FormatFloat(util.RoundFloat64(data.RequestAverageCostTime, 2), 'f', -1, 64) + "ms"
		requestCostMaxTime = strconv.FormatFloat(util.RoundFloat64(data.RequestCostMaxTime, 2), 'f', -1, 64) + "ms"
		requestCostMinTime = strconv.FormatFloat(util.RoundFloat64(data.RequestCostMinTime, 2), 'f', -1, 64) + "ms"
		requestAverageResponseTime = strconv.FormatFloat(util.RoundFloat64(data.RequestAverageResponseTime, 2), 'f', -1, 64) + "ms"
		requestResponseMaxTime = strconv.FormatFloat(util.RoundFloat64(data.RequestResponseMaxTime, 2), 'f', -1, 64) + "ms"
		requestResponseMinTime = strconv.FormatFloat(util.RoundFloat64(data.RequestResponseMinTime, 2), 'f', -1, 64) + "ms"
		transactionRate = strconv.FormatFloat(util.RoundFloat64(data.TransactionRate, 2), 'f', -1, 64) + " req/s"
		succTransactions = strconv.Itoa(data.SuccTransactions)
		failTransactions = strconv.Itoa(data.FailTransactions)
		timeOutTransactions = strconv.Itoa(data.TimeOutTransactions)
		dataTransferred = util.ByteToMB(data.DataTransferred)
		throughput = util.ByteToMB(int(data.Throughput)) + "/s"
		totalCostTime = strconv.FormatFloat(util.RoundFloat64(data.RequestCostMaxTime, 2), 'f', -1, 64) + "ms"
	}

	res := GetTaskDetailResponse{
		TargetAddress:              data.TargetAddress,
		TargetPort:                 data.TargetPort,
		Timeout:                    strconv.Itoa(data.Timeout) + "ms",
		ReadTimeout:                readTimeout,
		ExpectedBytes:              expectedBytes,
		ThreadNum:                  strconv.Itoa(data.ThreadNum),
		IsRepeat:                   data.IsRepeat,
		RepeatTime:                 repeatTime,
		IntervalTime:               sendInterval,
		HasResponse:                data.HasResponse,
		DataMapArr:                 nil,
		Taskid:                     data.Taskid,
		State:                      state,
		StartTime:                  startTime,
		EndTime:                    endTime,
		TotalRequestCount:          totalRequestCount,
		RequestAverageCostTime:     requestAverageCostTime,
		RequestCostMaxTime:         requestCostMaxTime,
		RequestCostMinTime:         requestCostMinTime,
		RequestAverageResponseTime: requestAverageResponseTime,
		RequestResponseMaxTime:     requestResponseMaxTime,
		RequestResponseMinTime:     requestResponseMinTime,
		TransactionRate:            transactionRate,
		SuccTransactions:           succTransactions,
		FailTransactions:           failTransactions,
		TimeOutTransactions:        timeOutTransactions,
		DataTransferred:            dataTransferred,
		Throughput:                 throughput,
		TotalCostTime:              totalCostTime,
	}

	return res
}