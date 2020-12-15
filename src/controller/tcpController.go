package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"strconv"
	"stress-tool/comon/util"
	"stress-tool/model"
	"stress-tool/src/tcp"
)

// 向外导出的方法变量字典
var TcpControllerMethodHandleMap map[string]func([]byte)interface{}
var log *go_logger.Logger

func init()  {
	log = util.GetLogger()
	// 容量会自增长
	TcpControllerMethodHandleMap = make(map[string]func([]byte)interface{}, 0)
	// 增加的方法要注册到 TcpControllerMethodHandleMap 列表中
	TcpControllerMethodHandleMap["/TestConnectivity"] = TestConnectivity
	TcpControllerMethodHandleMap["/CreateTask"] = CreateTask
	TcpControllerMethodHandleMap["/GetAllTaskDesc"] = GetAllTaskDesc
}

func TestConnectivity(request []byte) interface{} {
	var req model.TestConnectivityRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	m := tcp.TestConnectivityCheckParam(req)
	if m != "" {
		return util.ResponsePack(m, "", nil)
	}

	result := tcp.TestConnectivity(req.TargetAddress, req.TargetPort)
	var msg string
	if (result) {
		msg = "连接测试成功"
	} else {
		msg = "连接失败，请检查地址和端口"
	}
	return util.ResponseSuccPack(gin.H{
		"result": result,
		"msg": msg,
	})
}

/**
 * model.CreateTaskData 返回值
 * resCode string 错误码，为空则没错误
 */
func checkToCreateTaskData(req model.CreateTaskRequest) (res model.CreateTaskData, resCode string) {
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

	dataMapArr := req.DataMapArr
	var databytes []byte
	for _, v := range dataMapArr {
		t := v.Type
		lengthStr := v.Length
		data := v.Data
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			resCode += "1020"
			return
		}
		switch t {
		case model.NUMBER:
			switch length {
			case 1:
				dataInt, err := strconv.ParseInt(data, 10, 8)
				if err != nil {
					resCode += "1030"
					return
				}
				databytes = append(databytes, byte(int8(dataInt)))
			case 2:
				dataInt, err := strconv.ParseInt(data, 10, 16)
				if err != nil {
					resCode += "1031"
					return
				}
				databytes = append(databytes, byte(int16(dataInt)))
			case 4:
				dataInt, err := strconv.ParseInt(data, 10, 32)
				if err != nil {
					resCode += "1032"
					return
				}
				databytes = append(databytes, byte(int32(dataInt)))
			case 8:
				dataInt, err := strconv.ParseInt(data, 10, 64)
				if err != nil {
					resCode += "1033"
					return
				}
				databytes = append(databytes, byte(int64(dataInt)))
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
				databytes = append(databytes, byte(float32(dataFloat)))
			case 8:
				dataFloat, err := strconv.ParseFloat(data, 64)
				if err != nil {
					resCode += "1051"
					return
				}
				databytes = append(databytes, byte(dataFloat))
			default:
				resCode += "1060"
			}
		case model.STRING:
			length = len(data)
			bytes := []byte(data)
			databytes = append(databytes, bytes...)
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
		Data:          databytes,
	}

	return
}
func CreateTask(request []byte) interface{} {
	var req model.CreateTaskRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	res, resCode := checkToCreateTaskData(req)
	if resCode != "" {
		return util.ResponsePack(resCode, "", nil)
	}
	_, err = tcp.CreateTask(res)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	return util.ResponseSuccPack("创建任务成功")
}

func GetAllTaskDesc(request []byte) interface{}  {
	var req model.GetAllTaskDescRequest
	req.State = -1 // 默认 -1
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	allTaskDescript := tcp.GetAllTaskDescript(req.State)
	return util.ResponseSuccPack(allTaskDescript)
}












