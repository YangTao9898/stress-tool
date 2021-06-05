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
	"strings"
	"sync"
	"time"
)

// 向外导出的方法变量字典
var TcpControllerMethodHandleMap map[string]func([]byte, *gin.Context) interface{}
var log *go_logger.Logger
var tcpTaskFileLock sync.Mutex
var tcpTaskFileMap map[string]*model.SaveTcpTaskFileItem
var tcpTaskFileArr []*model.SaveTcpTaskFileItem

var tcpTestReturnConnMap map[string]*model.TcpReturnTestConnStruct
var tcpTestReturnConnMapLock sync.Mutex

const (
	tcpTaskFilePath                          = "./web-template/save/tcpTask.txt"
	tcpTaskTmpFilePath                       = "./web-template/save/.tcpTask.tmp.txt"
	TcpTestReturnSendRequestWaitResponseTime = 3000 // ms
)

func init() {
	log = util.GetLogger()
	// 容量会自增长
	TcpControllerMethodHandleMap = make(map[string]func([]byte, *gin.Context) interface{}, 0)
	tcpTestReturnConnMap = make(map[string]*model.TcpReturnTestConnStruct, 4)
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
	TcpControllerMethodHandleMap["/TcpTestReturnDeleteRequestQueue"] = TcpTestReturnDeleteRequestQueue
	TcpControllerMethodHandleMap["/TcpTestReturnRequestQueueUpdateData"] = TcpTestReturnRequestQueueUpdateData
	TcpControllerMethodHandleMap["/TcpTestReturnGetRequestQueue"] = TcpTestReturnGetRequestQueue
	TcpControllerMethodHandleMap["/TcpTestReturnSendRequest"] = TcpTestReturnSendRequest
	TcpControllerMethodHandleMap["/BinaryConvert"] = BinaryConvert
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

func createSession(connStruct *model.TcpReturnTestConnStruct, c *gin.Context) error {
	uniqueKey, err := util.GetUniqueString("TcpTestReturnConnect-")
	if err != nil {
		goto err
	}
	err = util.SetSession(c, sessionName, uniqueKey)
	if err != nil {
		goto err
	}
	tcpTestReturnConnMapLock.Lock()
	tcpTestReturnConnMap[uniqueKey] = connStruct
	tcpTestReturnConnMapLock.Unlock()
	return nil
err:
	log.Error("创建session发生错误，如存在连接则关闭连接：" + err.Error())
	if connStruct.IsConn {
		err := connStruct.Conn.Close()
		if err != nil {
			log.Error("关闭连接失败：" + err.Error())
		} else {
			log.Info("关闭连接成功")
		}
	}
	return util.WrapError("创建session发生错误", err)
}

func createConn(connStruct *model.TcpReturnTestConnStruct, targetAddress, targetPort string) error {
	conn, err := tcp.TcpConn(targetAddress, targetPort)
	if err != nil {
		return util.WrapError("连接失败", err)
	}
	if connStruct == nil {
		return util.NewErrorf("参数 connStruct 不能为空")
	}
	connStruct.Conn = conn
	connStruct.IsConn = true
	connStruct.Address = targetAddress
	connStruct.Port = targetPort
	return nil
}

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
	var connStruct *model.TcpReturnTestConnStruct
	if value == nil {
		goto END
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
	if connStruct == nil {
		connStruct = &model.TcpReturnTestConnStruct{}
	}
	err = createConn(connStruct, req.TargetAddress, req.TargetPort)
	if err != nil {
		ret.Msg = "创建连接发生错误：" + err.Error()
		log.Error(ret.Msg)
		return util.ResponseFailPack(ret.Msg)
	}

	err = createSession(connStruct, c)
	if err != nil {
		ret.Msg = "创建连接发生错误：" + err.Error()
		log.Error(ret.Msg)
		return util.ResponseSuccPack(ret)
	}

	ret.Msg = "连接成功"
	ret.Result = true
	return util.ResponseSuccPack(ret)
}

func getTcpTestReturnSession(c *gin.Context) (*model.TcpReturnTestConnStruct, string, error) {
	value := util.GetSession(c, sessionName)
	if value == nil {
		return nil, "", nil
	}

	var errMsg string
	uniqueKey, ok := value.(string)
	if !ok {
		errMsg = fmt.Sprintf("获取 session 发生错误：value[%+v]类型转换异常", value)
		return nil, "", util.NewErrorf(errMsg)
	}

	connStruct, ok := tcpTestReturnConnMap[uniqueKey]
	if !ok {
		errMsg = "该 session 数据不存在"
		return nil, "", util.NewErrorf(errMsg)
	}
	return connStruct, uniqueKey, nil
}

func TcpTestReturnDisconnect(request []byte, c *gin.Context) interface{} {
	var errMsg string
	connStruct, _, err := getTcpTestReturnSession(c)
	if connStruct == nil {
		errMsg = "断开连接发生错误：session 不存在"
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	if err != nil {
		errMsg = "断开连接发生错误：" + err.Error()
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}

	if connStruct.Conn == nil {
		errMsg = "断开连接发生错误：连接为空"
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	err = connStruct.Conn.Close()
	if err != nil { // 该错误不影响断开连接
		errMsg = "断开连接发生错误：" + err.Error()
		log.Error(errMsg)
	}
	/*err = util.DeleteSession(c, sessionName)
	if err != nil {
		errMsg = fmt.Sprintf("断开连接发生错误：%s", err.Error())
		log.Error(errMsg)
		return util.ResponseSuccPack(errMsg)
	}*/
	// 修改状态
	connStruct.IsConn = false
	connStruct.Conn = nil

	return util.ResponseSuccPack("断开连接成功")
}

// 删除请求数据队列，保存到 session
func TcpTestReturnDeleteRequestQueue(request []byte, c *gin.Context) interface{} {
	var errMsg string
	connStruct, _, err := getTcpTestReturnSession(c)
	if connStruct == nil {
		log.Debugf("session 不存在，无法删除请求数据队列")
		return util.ResponseSuccPack(nil)
	}
	if err != nil {
		errMsg = "删除请求队列发生错误：" + err.Error()
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	mapArrList := connStruct.DataMapArrList

	var req model.NumberRequest
	err = json.Unmarshal(request, &req)
	if err != nil {
		errMsg = "删除请求队列发生错误：" + err.Error()
		log.Error(errMsg)
		return util.ResponseSuccPack(errMsg)
	}

	if mapArrList == nil || len(mapArrList) == 0 {
		errMsg = "删除请求队列发生错误：session 数据为空，删除失败"
		log.Error(errMsg)
		return util.ResponseSuccPack(errMsg)
	}
	// 删除指定索引的元素
	connStruct.DataMapArrList = append(mapArrList[:req.Num], mapArrList[req.Num+1:]...)

	return util.ResponseSuccPack("删除请求队列成功")
}

// 更新请求数据队列的数据，保存到 session
func TcpTestReturnRequestQueueUpdateData(request []byte, c *gin.Context) interface{} {
	var errMsg string
	connStruct, _, err := getTcpTestReturnSession(c)
	if connStruct == nil {
		log.Infof("session 为空，创建 session, ip: [%s]", c.Request.Host)
		err = createSession(&model.TcpReturnTestConnStruct{}, c)
		if err != nil {
			errMsg = "创建连接发生错误：" + err.Error()
			log.Error(errMsg)
			return util.ResponseFailPack(errMsg)
		}
		connStruct, _, err = getTcpTestReturnSession(c)
	}
	if err != nil {
		errMsg = "添加请求队列发生错误：" + err.Error()
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	if connStruct == nil {
		errMsg = "session信息不存在，请重新连接"
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}
	mapArrList := connStruct.DataMapArrList

	var req model.RequestQueueUpdateDataRequest
	err = json.Unmarshal(request, &req)
	if err != nil {
		errMsg = "创建请求队列发生错误：" + err.Error()
		log.Error(errMsg)
		return util.ResponseSuccPack(errMsg)
	}

	if len(mapArrList) <= req.QueueIndex {
		connStruct.DataMapArrList = append(mapArrList, req.DataMapArr)
	} else {
		mapArrList[req.QueueIndex] = req.DataMapArr
	}
	return util.ResponseSuccPack("更新请求数据队列的数据成功")
}

// 获取 session 中的请求数据队列
func TcpTestReturnGetRequestQueue(request []byte, c *gin.Context) interface{} {
	connStruct, _, err := getTcpTestReturnSession(c)
	if err != nil {
		log.Debugf("session 获取失败，不填充请求队列：" + err.Error())
		return util.ResponseSuccPack(nil)
	}
	if connStruct == nil {
		log.Debugf("session 不存在，不填充请求队列")
		return util.ResponseSuccPack(nil)
	}
	address := connStruct.Address
	port := connStruct.Port
	res := model.TcpTestReturnGetRequestQueueResponse{
		Address:        address,
		Port:           port,
		IsConn:         connStruct.IsConn,
		DataMapArrList: connStruct.DataMapArrList,
	}

	return util.ResponseSuccPack(res)
}

// 测试返回界面发送请求到目标主机
func TcpTestReturnSendRequest(request []byte, c *gin.Context) interface{} {
	var errMsg string
	var req model.TcpTestReturnSendRequestRequest
	err := json.Unmarshal(request, &req)
	if err != nil {
		errMsg = util.GetErrorStack("参数转换错误", err)
		log.Error(errMsg)
		return util.ResponseSuccPack(errMsg)
	}
	m := tcp.TcpConnRequestCheckParam(req.TcpConnRequest)
	if m != "" {
		return util.ResponsePack(m, "", nil)
	}

	connStruct, _, err := getTcpTestReturnSession(c)
	if err != nil || connStruct == nil || !connStruct.IsConn {
		errMsg = "没有 session 连接信息，将建立连接"
		log.Info(errMsg)

		if connStruct == nil {
			connStruct = &model.TcpReturnTestConnStruct{}
		}

		err = createConn(connStruct, req.TargetAddress, req.TargetPort)
		if err != nil {
			errMsg = util.GetErrorStack("创建连接发生错误", err)
			log.Error(errMsg)
			return util.ResponseFailPack(errMsg)
		}

		err = createSession(connStruct, c)
		if err != nil {
			errMsg = util.GetErrorStack("创建连接发生错误", err)
			log.Error(errMsg)
			return util.ResponseSuccPack(errMsg)
		}
	}

	bytes, resCode, _, _, err := convert.DataMapToBytes(req.DataMapArr)
	if err != nil {
		return util.ResponsePack(resCode, err.Error(), nil)
	}

	_, err = connStruct.Conn.Write(bytes)
	if err != nil {
		connStruct.Conn.Close()
		connStruct.IsConn = false
		return util.ResponseFailPack("发送数据失败，请重新建立连接再试：" + err.Error())
	}

	err = connStruct.Conn.SetReadDeadline(time.Now().Add(time.Duration(TcpTestReturnSendRequestWaitResponseTime * 1e6)))
	if err != nil {
		connStruct.Conn.Close()
		connStruct.IsConn = false
		return util.ResponseFailPack("设置读取超时时间失败：" + err.Error())
	}

	res := model.TcpTestReturnSendRequestResponse{}
	connStruct.Bytes = nil
	for {
		buf := make([]byte, 1024)
		rdn, err := connStruct.Conn.Read(buf)
		connStruct.Bytes = append(connStruct.Bytes, buf[:rdn]...)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				// 根据读超时时间自动断开的连接，不算错误
			} else {
				res.ErrMsg = err.Error()
			}
			break
		}
	}

	// 将 byte 发送到前端界面
	if len(connStruct.Bytes) > 0 {
		res.BinaryStr = util.BytesToBinaryString(connStruct.Bytes, true)
	} else {
		return util.ResponseSuccPack(gin.H{
			"errMsg": "没有收到服务端发来的任何消息",
		})
	}

	return util.ResponseSuccPack(res)
}

func BinaryConvert(request []byte, c *gin.Context) interface{} {
	var req model.TcpTestReturnConvertRequest
	var errMsg string
	var retStr string

	err := json.Unmarshal(request, &req)
	if err != nil {
		errMsg = util.GetErrorStack("参数转换错误", err)
		log.Error(errMsg)
		return util.ResponseFailPack(errMsg)
	}

	if strings.TrimSpace(req.BinaryString) == "" {
		return util.ResponseFailPack("返回数据为空，不能够解析数据")
	}

	bytes, err := util.BinaryStringToBytes(req.BinaryString)
	if err != nil {
		return util.ResponseFailPack("解析二进制数据错误, 原因：" + err.Error())
	}
	var index int
	length := len(bytes)
	finish := false
	for _, v := range req.ConvertTypeArr {
		// 长度校验，int和float需要校验
		next := index + v.Length
		if v.Length > length {
			next = length
			finish = true
			switch v.Type {
			case "1":
				fallthrough
			case "2":
				if v.Length != 1 && v.Length != 2 && v.Length != 4 && v.Length != 8 {
					return util.ResponseFailPack(fmt.Sprintf("int类型不支持的长度[%d]", v.Length))
				}
			case "3":
				fallthrough
			case "4":
				if v.Length != 4 && v.Length != 8 {
					return util.ResponseFailPack(fmt.Sprintf("float类型不支持的长度[%d]", v.Length))
				}
			}
		}
		str, err := convert.ConvertBytesByType(bytes[index:next], v.Type)
		if err != nil {
			return util.ResponseFailPack("转换二进制数据失败，错误原因：" + err.Error())
		}
		index = next
		retStr += str
		if finish {
			break
		}
	}

	return util.ResponseSuccPack(retStr)
}
