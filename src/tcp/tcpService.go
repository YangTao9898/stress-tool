package tcp

import (
	"encoding/json"
	"errors"
	go_logger "github.com/phachon/go-logger"
	"net"
	"sort"
	"strconv"
	"stress-tool/comon/util"
	"stress-tool/model"
	"strings"
	"sync"
	"time"
)

const default_read_buf_size int = 1024

var taskMap map[string]*model.TaskDealData
var dataIntLength map[int]bool
var dataFloatLength map[int]bool
var log *go_logger.Logger

func init() {
	taskMap = make(map[string]*model.TaskDealData)
	dataIntLength = map[int]bool {
		1: true,
		2: true,
		4: true,
		8: true,
	}
	dataFloatLength = map[int]bool {
		4: true,
		8: true,
	}
}


// 检查数据范围是否连续
// return: int 为范围最大值 string 为返回的错误信息
func CheckDataRangeIsContinuous(mArr []map[string]int) (int, string) {
	rangelistMap  := make(map[int][2]int)
	var l []int = make([]int, len(mArr))
	index := 0

	m := make(map[string]int, len(mArr))
	// 合并多个Map
	for _, v := range mArr {
		for k, o := range v {
			m[k] = o
		}
	}
	for key, value := range m {
		nums := strings.Split(key, "~")
		if len(nums) != 2{
			return -1, "数据范围须遵循 0~4 格式"
		}
		var arr [2]int
		var err1, err2 error
		arr[0], err1 = strconv.Atoi(nums[0])
		arr[1], err2 = strconv.Atoi(nums[1])
		if err1 != nil || err2 != nil {
			return -1, "数据范围须遵循 0~4 格式"
		}
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
		l[index] = arr[0]
		index++
	}
	// 测试 key 对应的 value 所表示的范围是否连续
	sort.Ints(l)
	if l[0] != 0 { // 下标不是从0开始的
		return -1, "数据范围不是从0开始"
	}
	var arr [2]int
	for index, value := range l {
		tempArr := rangelistMap[value]
		if (index == 0) {
			arr = tempArr
			continue
		}
		if arr[1] + 1 != tempArr[0] {
			return -1, "数据范围不连续"
		}
		arr = tempArr
	}
	// 取出索引范围内最大的值, 加 1 得到数据长度
	return rangelistMap[l[len(l) - 1]][1] + 1, ""
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

	length, msg := CheckDataRangeIsContinuous(data.DataTypeMap)
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

// 创建测试任务
func CreateTask(req model.CreateTaskData) (string, error) {
	msg := checkCreateTaskData(req)
	if msg != "" {
		return "", errors.New(msg)
	}

	taskid := util.GetNowDateToStr(util.TIME_TEMPLATE_2);

	dealData := model.TaskDealData {
		CreateTaskData: req,
		Taskid: taskid,
		State: model.NOT_START,
	}
	// 未作键值 taskid 重复校验，因为单人使用并发量不大
	taskMap[taskid] = &dealData
	return taskid, nil
}

/**
 * 查询指定 state 状态的任务，如果 state 为-1则查询所有的任务
 */
func GetAllTaskDescript(state int) []model.TaskDealDataDescript {
	var arr []model.TaskDealDataDescript
	for _, v := range taskMap {
		if (state == -1 || state == v.State) {
			desc := model.TaskDealDataDescript{
				TaskId:       v.Taskid,
				TargetAddress: v.TargetAddress,
				TargetPort:   v.TargetPort,
				State:        v.State,
				StartTime:    v.StartTime,
				EndTime:      v.EndTime,
			}
			arr = append(arr, desc)
		}
	}
	return arr
}

func GetTaskByTaskId(taskId string) *model.TaskDealData {
	return taskMap[taskId]
}

func TestConnectivityCheckParam(req model.TestConnectivityRequest) string {
	resultCodes := ""
	if (req.TargetAddress == "") {
		resultCodes += "1001,"
	}
	if (req.TargetPort == "") {
		resultCodes += "1002,"
	}
	return resultCodes
}

// 测试连接
func TestConnectivity(address, port string) bool {
	conn, err := net.Dial("tcp", address + ":" + port)
	if err == nil {
		defer conn.Close()
	}
	return err == nil// && panicError == nil
}

type countData struct {
	totalRequestCount int
	totalCostTime float64
	costMaxTime float64
	costMinTime float64
	totalResponseTime float64
	responseMaxTime float64
	responseMinTime float64
	succCount int
	failCount int
	timeoutCount int
	sendDataSize int // 发送的数据量，单位 byte
	reciveDataSize int
}

// 如果测试失败了，只看返回值 countData 的 failCount 和 timeoutCount 计数
// cdata 里面的 max 和 min time 都不设置，只设置总时间
// syncCh 不为空时，支持同步
// return: countData 为计数信息，bool 为 true 时，测试成功，false 测试失败
func StartSingleThread(d model.CreateTaskData, outdata *countData, outb *bool, syncCh chan <- int) {
	cdata := outdata
	*outb = false
	if syncCh != nil {
		defer func() {
			syncCh <- 0
		}()
	}
	var conn net.Conn
	var err error
	timeoutSecond := d.Timeout
	if timeoutSecond > 0 {
		conn, err = net.DialTimeout("tcp", d.TargetAddress + ":" + d.TargetPort, time.Duration(timeoutSecond * 1e6))
		//conn, err = net.Dial("tcp", d.TargetAddress + ":" + d.TargetPort)
	} else {
		conn, err = net.Dial("tcp", d.TargetAddress + ":" + d.TargetPort)
	}

	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			cdata.timeoutCount++
		} else {
			cdata.failCount++
		}
		log.Error("StartSingleThread Dial err:" + err.Error())
		return
	}
	defer conn.Close()

	waitResponseCh := make(chan int)
	var responseStart time.Time
	var responseEnd time.Time
	if (d.HasResponse) { // 读取数据
		go func() {
			count := 0
			defer close(waitResponseCh)
			responseStart = time.Now()
			conn.SetReadDeadline(responseStart.Add(time.Duration(d.ReadTimeout * 1e6)))
			var buf []byte = make([]byte, default_read_buf_size)
			var preResponseTime time.Time
			for { // 不计算最后一次超时时间
				preResponseTime = time.Now()
				rdn, err := conn.Read(buf)
				count++
				responseEnd = time.Now()
				// fmt.Println("StartSingleThread read:", rdn, "byte:", string(buf))
				buf = make([]byte, default_read_buf_size)
				if err != nil {
					if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
						// 根据读超时时间自动断开的连接，不算错误，不计算本次响应时间
						if (count > 0) { // 如果一次数据都没接收，则算上超时时间
							responseEnd = preResponseTime
						} else {
							log.Error("StartSingleThread Read timeout err:" + err.Error())
							cdata.timeoutCount++
							return
						}
						break
					}
					log.Error("StartSingleThread Read err:" + err.Error())
					return
				}
				cdata.reciveDataSize += rdn
				if d.ExpectedBytes > 0 && cdata.reciveDataSize >= d.ExpectedBytes {
					break
				}
			}
		}()
	}

	writeStart := time.Now()
	conn.SetWriteDeadline(writeStart.Add(time.Duration(d.ReadTimeout * 1e9)))
	wrn, err := conn.Write(d.Data)
	if err != nil {
		cdata.failCount++
		log.Error("StartSingleThread Write err:" + err.Error())
		return
	}
	writeEnd := time.Now()
	cdata.sendDataSize += wrn

	<-waitResponseCh
	writeTotal := float64(util.GetTimestamp(writeEnd) - util.GetTimestamp(writeStart))
	if (d.HasResponse) {
		responseTotal := float64(util.GetTimestamp(responseEnd) - util.GetTimestamp(responseStart))
		cdata.totalCostTime = writeTotal + responseTotal
		cdata.totalResponseTime = responseTotal
	} else {
		cdata.totalCostTime = writeTotal
	}
	cdata.succCount++

	*outb = true
}

func compute(cdata countData, totalCount *countData, isSucc bool) {
	totalCount.totalRequestCount++
	if isSucc {
		totalCount.succCount += cdata.succCount
		totalCount.totalResponseTime += cdata.totalResponseTime
		totalCount.totalCostTime += cdata.totalCostTime
		if totalCount.costMaxTime < cdata.totalCostTime {
			totalCount.costMaxTime = cdata.totalCostTime
		}
		if totalCount.costMinTime > cdata.totalCostTime {
			totalCount.costMinTime = cdata.totalCostTime
		}
		if totalCount.responseMaxTime < cdata.totalResponseTime {
			totalCount.responseMaxTime = cdata.totalResponseTime
		}
		if totalCount.responseMinTime > cdata.totalResponseTime {
			totalCount.responseMinTime = cdata.totalResponseTime
		}
		totalCount.sendDataSize += cdata.sendDataSize
		totalCount.reciveDataSize += cdata.reciveDataSize
	} else {
		totalCount.timeoutCount += cdata.timeoutCount
		totalCount.failCount += cdata.failCount
	}
}

func startTestAndCompute(c model.CreateTaskData, outTotalCData *countData, syncCh chan int, mutex *sync.Mutex) {
	ch := make(chan int)
	var tempCData countData
	var ret bool
	go StartSingleThread(c, &tempCData, &ret, ch)
	// 计算结果
	<-ch
	mutex.Lock()
	compute(tempCData, outTotalCData, ret)
	mutex.Unlock()
	syncCh <- 0
}

// 开始任务
func StartTask(taskId string) error {
	v, ok := taskMap[taskId]
	if !ok {
		return errors.New("ID为[" + taskId + "]的任务不存在")
	}
	if v.State != model.NOT_START {
		return errors.New("该任务正在运行")
	}

	startTime := time.Now()
	v.State = model.RUNNING
	v.StartTime = util.GetDateToStr(startTime, util.TIME_TEMPLATE_1)

	repeatCount := 1
	if v.CreateTaskData.IsRepeat {
		repeatCount = v.CreateTaskData.RepeatTime
	}

	totalCData := countData{
		totalRequestCount: 0,
		totalCostTime:     0,
		costMaxTime:       -999,
		costMinTime:       999,
		totalResponseTime: 0,
		responseMaxTime:   -999,
		responseMinTime:   999,
		succCount:         0,
		failCount:         0,
		timeoutCount:      0,
		sendDataSize:      0,
		reciveDataSize:    0,
	}
	var lock sync.Mutex
	ch := make(chan int)
	for i := 0; i < repeatCount; i++ {
		go startTestAndCompute(v.CreateTaskData, &totalCData, ch, &lock)
	}
	for i := 0; i < repeatCount; i++ {
		<-ch
	}
	// 计算各类平均值
	v.EndTime = util.GetNowDateToStr(util.TIME_TEMPLATE_1)
	v.TotalRequestCount = totalCData.totalRequestCount
	v.RequestAverageCostTime = totalCData.totalCostTime / float64(totalCData.succCount)
	v.RequestCostMaxTime = totalCData.costMaxTime
	v.RequestCostMinTime = totalCData.costMinTime
	v.RequestAverageResponseTime = totalCData.totalResponseTime / float64(totalCData.succCount)
	v.RequestResponseMaxTime = totalCData.responseMaxTime
	v.RequestResponseMinTime = totalCData.responseMinTime
	if totalCData.totalCostTime == 0 {
		v.TransactionRate = 0
		v.Throughput = 0
	} else {
		v.TransactionRate = float64(totalCData.totalRequestCount) / (totalCData.totalCostTime / 1000)
		v.Throughput = float64(totalCData.sendDataSize + totalCData.reciveDataSize) / (totalCData.totalCostTime / 1000)
	}
	v.SuccTransactions = totalCData.succCount
	v.FailTransactions = totalCData.failCount
	v.TimeOutTransactions = totalCData.timeoutCount
	v.DataTransferred = totalCData.reciveDataSize + totalCData.sendDataSize
	v.TotalCostTime = totalCData.totalCostTime
	v.State = model.FINISH

	return nil
}

// 将 TaskDealData 转为 json
func ConvertTaskDealDataToJson(d model.TaskDealData) (string, error) {
	bytes, e := json.Marshal(d)
	return string(bytes), e
}







