package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"stress-tool/comon/util"
	"stress-tool/model"
	"stress-tool/src/tcp"
)

// 向外导出的方法变量字典
var TcpControllerMethodHandleMap map[string]func([]byte) interface{}
var log *go_logger.Logger

func init() {
	log = util.GetLogger()
	// 容量会自增长
	TcpControllerMethodHandleMap = make(map[string]func([]byte) interface{}, 0)
	// 增加的方法要注册到 TcpControllerMethodHandleMap 列表中
	TcpControllerMethodHandleMap["/TestConnectivity"] = TestConnectivity
	TcpControllerMethodHandleMap["/CreateTask"] = CreateTask
	TcpControllerMethodHandleMap["/GetAllTaskDesc"] = GetAllTaskDesc
	TcpControllerMethodHandleMap["/GetTaskDetail"] = GetTaskDetail
	TcpControllerMethodHandleMap["/GetStrBytes"] = GetStrBytes
	TcpControllerMethodHandleMap["/StartTask"] = StartTask
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
	if result {
		msg = "连接测试成功"
	} else {
		msg = "连接失败，请检查地址和端口"
	}
	return util.ResponseSuccPack(gin.H{
		"result": result,
		"msg":    msg,
	})
}

func CreateTask(request []byte) interface{} {
	var req model.CreateTaskRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	res, resCode := model.CheckToCreateTaskData(req)
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

func GetAllTaskDesc(request []byte) interface{} {
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

func GetTaskDetail(request []byte) interface{} {
	var req model.TaskIdParamRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}

	taskDealData := tcp.GetTaskByTaskId(req.TaskId)
	if taskDealData == nil {
		return util.ResponseFailPack("该压测任务不存在")
	}
	response, err := model.TaskDealDataToGetTaskDetailResponse(*taskDealData)
	if err != nil {
		return util.ResponseFailPack(err.Error())
	}
	return util.ResponseSuccPack(response)
}

func GetStrBytes(request []byte) interface{} {
	return len(request)
}

func StartTask(request []byte) interface{} {
	var req model.TaskIdParamRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	err = tcp.StartTask(req.TaskId)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	return util.ResponsePack(util.RESULT_OK, "任务开始执行...", nil)
}
