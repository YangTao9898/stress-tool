package convert

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	go_logger "github.com/phachon/go-logger"
	"strconv"
	"stress-tool/comon/util"
	"stress-tool/model"
)

var log *go_logger.Logger

func init() {
	log = util.GetLogger()
}

/**
 * CreateTaskData 返回值
 * resCode string 错误码，为空则没错误
 */
func CheckToCreateTaskData(req model.CreateTaskRequest) (res model.CreateTaskData, resCode string) {
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

	readTimeout := 10000
	if req.HasResponse && req.ReadTimeout != "" {
		readTimeout, err = strconv.Atoi(req.ReadTimeout)
		if err != nil {
			resCode += "1005,"
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
		if req.IntervalTime != "" {
			intervalTime, err = strconv.Atoi(req.IntervalTime)
			if err != nil {
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
	isBigEndArr := make([]bool, len(dataMapArr))
	num := 0
	for index, v := range dataMapArr {
		lengthStr := v.Length
		data := v.Data
		isBigEnd := v.IsBigEnd
		isBigEndArr[index] = isBigEnd
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
		case model.NUMBER:
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
					err = binary.Write(databytes, binary.LittleEndian, int8(dataInt))
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
					err = binary.Write(databytes, binary.LittleEndian, int16(dataInt))
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
					err = binary.Write(databytes, binary.LittleEndian, int32(dataInt))
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
					err = binary.Write(databytes, binary.LittleEndian, dataInt)
				}
				if err != nil {
					log.Error(err.Error())
				}
			default:
				resCode += "1040"
				return
			}

		case model.FLOAT:
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
					err = binary.Write(databytes, binary.LittleEndian, float32(dataFloat))
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
					err = binary.Write(databytes, binary.LittleEndian, dataFloat)
				}
				if err != nil {
					log.Error(err.Error())
				}
			default:
				resCode += "1060"
			}
		case model.STRING:
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

	res = model.CreateTaskData{
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
		IsBigEnd:      isBigEndArr,
		Data:          databytes.Bytes(),
	}

	return
}

func byteToDataString(bs []byte, t int, length int, isBigEnd bool) (string, error) {
	buf := bytes.NewBuffer(bs)
	var resStr string
	switch t {
	case model.NUMBER:
		var err error
		if isBigEnd {
			switch length {
			case 1:
				var n int8
				err = binary.Read(buf, binary.BigEndian, &n)
				resStr = strconv.FormatInt(int64(n), 10)
			case 2:
				var n int16
				err = binary.Read(buf, binary.BigEndian, &n)
				resStr = strconv.FormatInt(int64(n), 10)
			case 4:
				var n int32
				err = binary.Read(buf, binary.BigEndian, &n)
				resStr = strconv.FormatInt(int64(n), 10)
			case 8:
				var n int64
				err = binary.Read(buf, binary.BigEndian, &n)
				resStr = strconv.FormatInt(n, 10)
			default:
				return "", util.NewErrorf("NUMBER not support data length[%d]", length)
			}
		} else {
			switch length {
			case 1:
				var n int8
				err = binary.Read(buf, binary.LittleEndian, &n)
				resStr = strconv.FormatInt(int64(n), 10)
			case 2:
				var n int16
				err = binary.Read(buf, binary.LittleEndian, &n)
				resStr = strconv.FormatInt(int64(n), 10)
			case 4:
				var n int32
				err = binary.Read(buf, binary.LittleEndian, &n)
				resStr = strconv.FormatInt(int64(n), 10)
			case 8:
				var n int64
				err = binary.Read(buf, binary.LittleEndian, &n)
				resStr = strconv.FormatInt(n, 10)
			default:
				return "", util.NewErrorf("NUMBER not support data length[%d]", length)
			}
		}
		if err != nil {
			return "", err
		}
		return resStr, err
	case model.FLOAT:
		var err error
		if isBigEnd {
			switch length {
			case 4:
				var f float32
				err = binary.Read(buf, binary.BigEndian, &f)
				resStr = strconv.FormatFloat(float64(f), 'f', -1, 32)
			case 8:
				var f float64
				err = binary.Read(buf, binary.BigEndian, &f)
				resStr = strconv.FormatFloat(f, 'f', -1, 64)
			default:
				return "", util.NewErrorf("FLOAT not support data length[%d]", length)
			}
		} else {
			switch length {
			case 4:
				var f float32
				err = binary.Read(buf, binary.LittleEndian, &f)
				resStr = strconv.FormatFloat(float64(f), 'f', -1, 32)
			case 8:
				var f float64
				err = binary.Read(buf, binary.LittleEndian, &f)
				resStr = strconv.FormatFloat(f, 'f', -1, 64)
			default:
				return "", util.NewErrorf("FLOAT not support data length[%d]", length)
			}
		}
		if err != nil {
			return "", err
		}
		return resStr, err
	case model.STRING:
		return string(bs), nil
	default:
		return "", util.NewErrorf("not support data type[%d]", t)
	}
}

func TaskDealDataToGetTaskDetailResponse(data model.TaskDealData) (model.GetTaskDetailResponse, error) {
	var res model.GetTaskDetailResponse
	readTimeout := "-"
	expectedBytes := "-"
	if data.HasResponse {
		readTimeout = strconv.Itoa(data.ReadTimeout) + " ms"
		if data.ExpectedBytes > 0 {
			expectedBytes = strconv.Itoa(data.ExpectedBytes) + " Byte"
		}
	}

	repeatTime := "-"
	sendInterval := "-"
	if data.IsRepeat {
		repeatTime = strconv.Itoa(data.RepeatTime)
		if data.IntervalTime > 0 {
			sendInterval = strconv.Itoa(data.IntervalTime) + " ms"
		}
	}

	state := ""
	switch data.State {
	case model.NOT_START:
		state = "未执行"
	case model.READY:
		state = "等待执行"
	case model.RUNNING:
		state = "正在执行"
	case model.FINISH:
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
	/*requestAverageResponseTime := "-"
	requestResponseMaxTime := "-"
	requestResponseMinTime := "-"*/
	threadAverageCostTime := "-"
	threadCostMaxTime := "-"
	threadCostMinTIme := "-"
	transactionRate := "-"
	succTransactions := "-"
	failTransactions := "-"
	timeOutTransactions := "-"
	dataTransferred := "-"
	throughput := "-"
	recvBytes := "-"
	totalCostTime := "-"
	totalRealCostTime := "-"
	if data.State == model.RUNNING {
		startTime = data.StartTime
	} else if data.State == model.FINISH {
		startTime = data.StartTime
		endTime = data.EndTime
		totalRequestCount = strconv.Itoa(data.TotalRequestCount)
		succTransactions = strconv.Itoa(data.SuccTransactions)
		failTransactions = strconv.Itoa(data.FailTransactions)
		timeOutTransactions = strconv.Itoa(data.TimeOutTransactions)
		dataTransferred = util.ByteToMB(data.DataTransferred)
		throughput = util.ByteToMB(int(data.Throughput)) + "/s"
		recvBytes = util.ByteToMB(data.RecvBytes)
		totalCostTime = strconv.FormatFloat(util.RoundFloat64(data.TotalCostTime, 2), 'f', -1, 64) + " ms"
		totalRealCostTime = strconv.FormatFloat(util.RoundFloat64(data.TotalRealCostTime, 2), 'f', -1, 64) + " ms"
		if data.SuccTransactions != 0 {
			requestAverageCostTime = strconv.FormatFloat(util.RoundFloat64(data.RequestAverageCostTime, 2), 'f', -1, 64) + " ms"
			requestCostMaxTime = strconv.FormatFloat(util.RoundFloat64(data.RequestCostMaxTime, 2), 'f', -1, 64) + " ms"
			requestCostMinTime = strconv.FormatFloat(util.RoundFloat64(data.RequestCostMinTime, 2), 'f', -1, 64) + " ms"
			/*requestAverageResponseTime = strconv.FormatFloat(util.RoundFloat64(data.RequestAverageResponseTime, 2), 'f', -1, 64) + " ms"
			requestResponseMaxTime = strconv.FormatFloat(util.RoundFloat64(data.RequestResponseMaxTime, 2), 'f', -1, 64) + " ms"
			requestResponseMinTime = strconv.FormatFloat(util.RoundFloat64(data.RequestResponseMinTime, 2), 'f', -1, 64) + " ms"*/
			transactionRate = strconv.FormatFloat(util.RoundFloat64(data.TransactionRate, 2), 'f', -1, 64) + " req/s"
		}
		threadAverageCostTime = strconv.FormatFloat(util.RoundFloat64(data.ThreadAverageCostTime, 2), 'f', -1, 64) + " ms"
		threadCostMaxTime = strconv.FormatInt(data.ThreadCostMaxTime, 10) + " ms"
		threadCostMinTIme = strconv.FormatInt(data.ThreadCostMinTime, 10) + " ms"
	}

	dataMapArr := make([]model.GetTaskDetailDataMap, len(data.DataTypeMap))
	for index, obj := range data.DataTypeMap {
		for k, dataType := range obj { // 仅一次
			n1, n2, err := model.CreateTaskDataKeySplite(k)
			if err != nil {
				errMsg := fmt.Sprintf("%s occur error: %s", "TaskDealDataToGetTaskDetailResponse", err.Error())
				return res, errors.New(errMsg)
			}
			length := n2 - n1 + 1
			if length <= 0 {
				return res, errors.New(fmt.Sprintf("[%s] connot less [%s]", n2, n1))
			}
			dataBytes := data.Data[n1 : n2+1]
			dataMapArr[index].Type = strconv.Itoa(dataType)
			dataMapArr[index].Length = strconv.Itoa(length)
			dataMapArr[index].BinaryData = util.BytesToBinaryString(dataBytes)
			dataMapArr[index].ByteData = util.BytesToByteString(dataBytes)
			dataMapArr[index].Data, err = byteToDataString(dataBytes, dataType, length, data.IsBigEnd[index])
			dataMapArr[index].IsBigEnd = data.IsBigEnd[index]
			if err != nil {
				return res, util.WrapError("TaskDealDataToGetTaskDetailResponse parse data err:", err)
			}
		}
	}

	timeout := "-"
	if data.Timeout > 0 {
		timeout = strconv.Itoa(data.Timeout) + " ms"
	}
	res = model.GetTaskDetailResponse{
		TargetAddress:          data.TargetAddress,
		TargetPort:             data.TargetPort,
		Timeout:                timeout,
		ReadTimeout:            readTimeout,
		ExpectedBytes:          expectedBytes,
		ThreadNum:              strconv.Itoa(data.ThreadNum),
		IsRepeat:               data.IsRepeat,
		RepeatTime:             repeatTime,
		IntervalTime:           sendInterval,
		HasResponse:            data.HasResponse,
		DataMapArr:             dataMapArr,
		Taskid:                 data.Taskid,
		State:                  state,
		StartTime:              startTime,
		EndTime:                endTime,
		TotalRequestCount:      totalRequestCount,
		RequestAverageCostTime: requestAverageCostTime,
		RequestCostMaxTime:     requestCostMaxTime,
		RequestCostMinTime:     requestCostMinTime,
		/*RequestAverageResponseTime: requestAverageResponseTime,
		RequestResponseMaxTime:     requestResponseMaxTime,
		RequestResponseMinTime:     requestResponseMinTime,*/
		ThreadAverageCostTime: threadAverageCostTime,
		ThreadCostMaxTime:     threadCostMaxTime,
		ThreadCostMinTime:     threadCostMinTIme,
		TransactionRate:       transactionRate,
		SuccTransactions:      succTransactions,
		FailTransactions:      failTransactions,
		TimeOutTransactions:   timeOutTransactions,
		DataTransferred:       dataTransferred,
		Throughput:            throughput,
		RecvBytes:             recvBytes,
		TotalCostTime:         totalCostTime,
		TotalRealCostTime:     totalRealCostTime,
	}

	return res, nil
}

func ToDescFromSaveTcpTaskFileItem(item []*model.SaveTcpTaskFileItem) []model.SaveTcpTaskFileDesc {
	var resArr []model.SaveTcpTaskFileDesc
	for i := len(item) - 1; i >= 0; i-- {
		var desc model.SaveTcpTaskFileDesc
		desc.SaveTaskId = (*item[i]).SaveTaskId
		desc.SaveTaskTag = (*item[i]).TaskData.SaveTaskTag
		desc.SaveTime = (*item[i]).SaveTime
		desc.TargetAddress = (*item[i]).TaskData.TargetAddress
		desc.TargetPort = (*item[i]).TaskData.TargetPort
		desc.ThreadNum = (*item[i]).TaskData.ThreadNum
		resArr = append(resArr, desc)
	}
	return resArr
}
