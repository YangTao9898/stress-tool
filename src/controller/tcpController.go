package controller

import (
	"bufio"
	"encoding/json"
	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"os"
	"stress-tool/comon/util"
	"stress-tool/model"
	"stress-tool/src/tcp"
	"sync"
	"time"
)

// 向外导出的方法变量字典
var TcpControllerMethodHandleMap map[string]func([]byte) interface{}
var log *go_logger.Logger
var tcpTaskFileLock sync.Mutex
var tcpTaskFileMap map[string]*model.SaveTcpTaskFileItem
var tcpTaskFileArr []*model.SaveTcpTaskFileItem

const (
	tcpTaskFilePath    = "./web-template/save/tcpTask.txt"
	tcpTaskTmpFilePath = "./web-template/save/.tcpTask.tmp.txt"
)

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
	TcpControllerMethodHandleMap["/SaveTask"] = SaveTask
	TcpControllerMethodHandleMap["/GetSaveTaskDesc"] = GetSaveTaskDesc
	TcpControllerMethodHandleMap["/GetSaveTaskDetail"] = GetSaveTaskDetail
	TcpControllerMethodHandleMap["/DeleteSaveTask"] = DeleteSaveTask
}

// 如未初始化相关变量则初始化
// 调用前请加锁
func initTcpTaskFile() error {
	if tcpTaskFileMap == nil {
		tcpTaskFileMap = make(map[string]*model.SaveTcpTaskFileItem)
		tcpTaskFileArr = make([]*model.SaveTcpTaskFileItem, 0)
		// 从文件中加载
		file, err := os.OpenFile(tcpTaskFilePath, os.O_APPEND|os.O_CREATE, 0644)
		defer file.Close()
		if err != nil {
			return util.WrapError("initTcpTaskFile err:", err)
		}
		// 逐行读取
		reader := bufio.NewReader(file)
		for {
			line, isPrefix, err := reader.ReadLine()
			if isPrefix {
				return util.NewErrorf("this line is too long, [%s]", line)
			}
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				return util.WrapError("initTcpTaskFile readLine err", err)
			}
			var o model.SaveTcpTaskFileItem
			err = json.Unmarshal(line, &o)
			if err != nil {
				return util.WrapError("initTcpTaskFile Unmarshal err", err)
			}
			tcpTaskFileMap[o.SaveTaskId] = &o
			tcpTaskFileArr = append(tcpTaskFileArr, &o)
		}
	}
	return nil
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

func SaveTask(request []byte) interface{} {
	var req model.CreateTaskRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	// 检查参数是否符合规范
	_, resCode := model.CheckToCreateTaskData(req)
	if resCode != "" {
		return util.ResponsePack(resCode, "", nil)
	}

	tcpTaskFileLock.Lock() // 加锁防止多线程写入
	defer tcpTaskFileLock.Unlock()
	err = initTcpTaskFile()
	if err != nil {
		log.Error(err.Error())
		return util.ResponseSuccPack("保存任务失败")
	}

	// 往文件追加内容
	file, err := os.OpenFile(tcpTaskFilePath, os.O_APPEND|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack("保存任务失败")
	}
	saveTime := time.Now()
	saveTaskId := util.GetDateToStrWithMillisecond(saveTime)
	saveTaskData := model.SaveTcpTaskFileItem{
		SaveTaskId: saveTaskId,
		SaveTime:   util.GetDateToStr(saveTime, util.TIME_TEMPLATE_1),
		TaskData:   req,
	}
	bytes, err := json.Marshal(saveTaskData)
	_, err = file.WriteString(string(bytes) + "\n")
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack("保存任务失败")
	}
	tcpTaskFileMap[saveTaskId] = &saveTaskData
	tcpTaskFileArr = append(tcpTaskFileArr, &saveTaskData)
	return util.ResponseSuccPack("保存任务成功")
}

func GetSaveTaskDesc(request []byte) interface{} {
	tcpTaskFileLock.Lock()
	defer tcpTaskFileLock.Unlock()
	err := initTcpTaskFile()
	if err != nil {
		log.Error(err.Error())
		return util.ResponseSuccPack("获取保存的任务失败")
	}
	resArr := model.ToDescFromSaveTcpTaskFileItem(tcpTaskFileArr)
	return util.ResponseSuccPack(resArr)
}

func GetSaveTaskDetail(request []byte) interface{} {
	tcpTaskFileLock.Lock()
	defer tcpTaskFileLock.Unlock()
	err := initTcpTaskFile()
	if err != nil {
		log.Error(err.Error())
		return util.ResponseSuccPack("加载保存的任务失败")
	}

	var req model.SaveTaskIdStruct
	err = json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseSuccPack("加载保存的任务失败")
	}
	if item, ok := tcpTaskFileMap[req.SaveTaskId]; ok {
		return util.ResponseSuccPack(item)
	} else {
		return util.ResponseFailPack("该保存的任务不存在")
	}
}

func DeleteSaveTask(request []byte) interface{} {
	var arr model.SaveTaskIdArrStruct
	err := json.Unmarshal(request, &arr)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseSuccPack("删除保任务失败")
	}
	if arr.SaveTaskIdArr == nil || len(arr.SaveTaskIdArr) == 0 {
		return util.ResponseFailPack("没有要删除的保存任务")
	}

	tcpTaskFileLock.Lock()
	defer tcpTaskFileLock.Unlock()
	// 删除 map 中对应的任务
	for _, v := range arr.SaveTaskIdArr {
		delete(tcpTaskFileMap, v)
	}

	succ := false
	defer func() {
		if !succ {
			tcpTaskFileMap = nil
		}
	}()
	// 将数据写入临时文件中
	file, err := os.OpenFile(tcpTaskTmpFilePath, os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		file.Close()
		log.Error(err.Error())
		return util.ResponseFailPack("删除保存任务失败")
	}

	var tempArr []*model.SaveTcpTaskFileItem
	// 删除 arr 中对应的任务
	for _, v := range tcpTaskFileArr {
		if _, ok := tcpTaskFileMap[v.SaveTaskId]; ok {
			tempArr = append(tempArr, v)
			bytes, err := json.Marshal(v)
			if err != nil {
				file.Close()
				log.Error(err.Error())
				return util.ResponseFailPack("删除保存任务失败")
			}
			_, err = file.WriteString(string(bytes) + "\n")
			if err != nil {
				file.Close()
				log.Error(err.Error())
				return util.ResponseFailPack("删除保存任务失败")
			}
		}
	}
	tcpTaskFileArr = tempArr

	file.Close()
	// 将原文件删除.用临时文件替换
	err = os.Remove(tcpTaskFilePath)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack("删除保存任务失败")
	}
	err = os.Rename(tcpTaskTmpFilePath, tcpTaskFilePath)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack("删除保存任务失败")
	}

	succ = true
	return util.ResponseSuccPack("删除保存任务成功")
}
