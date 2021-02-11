package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"net"
	"os"
	"stress-tool/comon/util"
	"stress-tool/model"
	"stress-tool/src/convert"
	"stress-tool/src/tcp"
	"sync"
	"time"
)

// 向外导出的方法变量字典
var TcpControllerMethodHandleMap map[string]func([]byte, *gin.Context) interface{}
var log *go_logger.Logger
var tcpTaskFileLock sync.Mutex
var tcpTaskFileMap map[string]*model.SaveTcpTaskFileItem
var tcpTaskFileArr []*model.SaveTcpTaskFileItem

var tcpTestReturnConnMap map[string]model.TcpReturnTestConnStruct

const (
	tcpTaskFilePath    = "./web-template/save/tcpTask.txt"
	tcpTaskTmpFilePath = "./web-template/save/.tcpTask.tmp.txt"
)

func init() {
	log = util.GetLogger()
	// 容量会自增长
	TcpControllerMethodHandleMap = make(map[string]func([]byte, *gin.Context) interface{}, 0)
	tcpTestReturnConnMap = make(map[string]model.TcpReturnTestConnStruct, 4)
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
	TcpControllerMethodHandleMap["/TcpTestReturnConnect"] = TcpTestReturnConnect
	TcpControllerMethodHandleMap["/TcpTestReturnDisconnect"] = TcpTestReturnDisconnect
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

func TestConnectivity(request []byte, c *gin.Context) interface{} {
	var req model.TcpConnRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	m := tcp.TcpConnRequestCheckParam(req)
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

func CreateTask(request []byte, c *gin.Context) interface{} {
	var req model.CreateTaskRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	res, resCode := convert.CheckToCreateTaskData(req)
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

func GetAllTaskDesc(request []byte, c *gin.Context) interface{} {
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

func GetTaskDetail(request []byte, c *gin.Context) interface{} {
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
	response, err := convert.TaskDealDataToGetTaskDetailResponse(*taskDealData)
	if err != nil {
		return util.ResponseFailPack(err.Error())
	}
	return util.ResponseSuccPack(response)
}

func GetStrBytes(request []byte, c *gin.Context) interface{} {
	return len(request)
}

func StartTask(request []byte, c *gin.Context) interface{} {
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

func SaveTask(request []byte, c *gin.Context) interface{} {
	var req model.CreateTaskRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack(err.Error())
	}
	// 检查参数是否符合规范
	_, resCode := convert.CheckToCreateTaskData(req)
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

func GetSaveTaskDesc(request []byte, c *gin.Context) interface{} {
	tcpTaskFileLock.Lock()
	defer tcpTaskFileLock.Unlock()
	err := initTcpTaskFile()
	if err != nil {
		log.Error(err.Error())
		return util.ResponseSuccPack("获取保存的任务失败")
	}
	resArr := convert.ToDescFromSaveTcpTaskFileItem(tcpTaskFileArr)
	return util.ResponseSuccPack(resArr)
}

func GetSaveTaskDetail(request []byte, c *gin.Context) interface{} {
	tcpTaskFileLock.Lock()
	defer tcpTaskFileLock.Unlock()
	err := initTcpTaskFile()
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack("加载保存的任务失败")
	}

	var req model.SaveTaskIdStruct
	err = json.Unmarshal(request, &req)
	if err != nil {
		log.Error(err.Error())
		return util.ResponseFailPack("加载保存的任务失败")
	}
	if item, ok := tcpTaskFileMap[req.SaveTaskId]; ok {
		return util.ResponseSuccPack(*item)
	} else {
		return util.ResponseFailPack("该保存的任务不存在")
	}
}

func DeleteSaveTask(request []byte, c *gin.Context) interface{} {
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

const sessionName = "TCP_RETURN_TEST_TCP_CONN_SESSION"

func TcpTestReturnConnect(request []byte, c *gin.Context) interface{} {
	ret := model.TcpTestReturnConnectResponse{
		Msg:    "",
		Result: false,
	}

	var req model.TcpConnRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		ret.Msg = "创建连接发生错误：" + err.Error()
		log.Error(ret.Msg)
		return util.ResponseSuccPack(ret)
	}
	m := tcp.TcpConnRequestCheckParam(req)
	if m != "" {
		return util.ResponsePack(m, "", nil)
	}

	// 查看连接是否存在
	value := util.GetSession(c, sessionName)
	var connStruct model.TcpReturnTestConnStruct
	if value == nil {
		connStruct = model.TcpReturnTestConnStruct{}
	} else {
		var ok bool
		var uniqueKey string
		uniqueKey, ok = value.(string)
		if !ok {
			ret.Msg = fmt.Sprintf("创建连接发生错误：value[%+v]类型转换异常", value)
			log.Error(ret.Msg)
			return util.ResponseSuccPack(ret)
		}
		connStruct, ok = tcpTestReturnConnMap[uniqueKey]
		if !ok || !connStruct.IsConn { // 不存在或者未连接
			goto END
		}
		ret.Msg = "连接已经存在"
		ret.Result = true
		return util.ResponseSuccPack(ret)
	}

END:
	conn, err := tcp.TcpConn(req.TargetAddress, req.TargetPort)
	if err != nil {
		ret.Msg = "创建连接发生错误：" + err.Error()
		log.Error(ret.Msg)
		return util.ResponseSuccPack(ret)
	}
	connStruct.Conn = conn
	connStruct.IsConn = true
	uniqueKey, err := util.GetUniqueString("TcpTestReturnConnect-")
	if err != nil {
		ret.Msg = "创建连接发生错误：" + err.Error()
		log.Error(ret.Msg)
		return util.ResponseSuccPack(ret)
	}
	err = util.SetSession(c, sessionName, uniqueKey)
	if err != nil {
		ret.Msg = "创建连接发生错误：" + err.Error()
		log.Error(ret.Msg)
		return util.ResponseSuccPack(ret)
	}
	tcpTestReturnConnMap[uniqueKey] = connStruct
	ret.Msg = "连接成功"
	ret.Result = true
	return util.ResponseSuccPack(ret)
}

func TcpTestReturnDisconnect(request []byte, c *gin.Context) interface{} {
	value := util.GetSession(c, sessionName)
	if value == nil {
		return util.ResponseFailPack("没有进行连接不能断开")
	}

	var errMsg string
	uniqueKey, ok := value.(string)
	if !ok {
		errMsg = fmt.Sprintf("断开连接发生错误：value[%+v]类型转换异常", value)
		log.Error(errMsg)
		return util.ResponseSuccPack(errMsg)
	}

	connStruct, ok := tcpTestReturnConnMap[uniqueKey]
	if !ok {
		errMsg = "断开连接发生错误：该连接不存在"
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	if connStruct.Conn == nil {
		errMsg = "断开连接发生错误：连接为空"
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	conn, ok := connStruct.Conn.(net.Conn)
	if !ok {
		errMsg = "断开连接发生错误：类型转换错误"
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	err := conn.Close()
	if err != nil { // 该错误不影响断开连接
		errMsg = "断开连接发生错误：" + err.Error()
		log.Error(errMsg)
	}
	err = util.DeleteSession(c, sessionName)
	if err != nil {
		errMsg = fmt.Sprintf("断开连接发生错误：%s", err.Error())
		log.Error(errMsg)
		return util.ResponseSuccPack(errMsg)
	}

	return util.ResponseSuccPack("断开连接成功")
}
