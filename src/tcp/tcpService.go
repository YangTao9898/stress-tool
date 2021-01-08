package tcp

import (
	"errors"
	go_logger "github.com/phachon/go-logger"
	"net"
	"stress-tool/comon/util"
	"stress-tool/model"
	"sync"
	"time"
)

const default_read_buf_size int = 1024

var taskMap map[string]*model.TaskDealData // 供查找用
var taskMapArr []*model.TaskDealData       // 供排序用
var dataIntLength map[int]bool
var dataFloatLength map[int]bool
var log *go_logger.Logger

func init() {
	log = util.GetLogger()
	taskMap = make(map[string]*model.TaskDealData)
	dataIntLength = map[int]bool{
		1: true,
		2: true,
		4: true,
		8: true,
	}
	dataFloatLength = map[int]bool{
		4: true,
		8: true,
	}
}

// 创建测试任务
func CreateTask(req model.CreateTaskData) (string, error) {
	msg := checkCreateTaskData(req)
	if msg != "" {
		return "", errors.New(msg)
	}

	taskid := util.GetDateToStrWithMillisecond(time.Now())

	dealData := model.TaskDealData{
		CreateTaskData: req,
		Taskid:         taskid,
		State:          model.NOT_START,
	}
	// 未作键值 taskid 重复校验，因为单人使用并发量不大
	taskMap[taskid] = &dealData
	taskMapArr = append(taskMapArr, &dealData)
	return taskid, nil
}

/**
 * 查询指定 state 状态的任务，如果 state 为-1则查询所有的任务
 */
func GetAllTaskDescript(state int) []model.TaskDealDataDescript {
	var arr []model.TaskDealDataDescript
	lenth := len(taskMapArr)
	for i := lenth - 1; i >= 0; i-- {
		v := taskMapArr[i]
		if state == -1 || state == v.State {
			desc := model.TaskDealDataDescript{
				TaskId:        v.Taskid,
				TargetAddress: v.TargetAddress,
				TargetPort:    v.TargetPort,
				State:         v.State,
				StartTime:     v.StartTime,
				EndTime:       v.EndTime,
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
	if req.TargetAddress == "" {
		resultCodes += "1001,"
	}
	if req.TargetPort == "" {
		resultCodes += "1002,"
	}
	return resultCodes
}

// 测试连接
func TestConnectivity(address, port string) bool {
	conn, err := net.Dial("tcp", address+":"+port)
	if err == nil {
		defer conn.Close()
	}
	return err == nil // && panicError == nil
}

// 如果测试失败了，只看返回值 countData 的 failCount 和 timeoutCount 计数
// cdata 里面的 max 和 min time 都不设置，只设置总时间
// syncCh 不为空时，支持同步
// return: countData 为计数信息，bool 为 true 时，测试成功，false 测试失败
func startSingleThread(d model.CreateTaskData, outdata *countData, outb *bool, syncCh chan<- int) {
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
		conn, err = net.DialTimeout("tcp", d.TargetAddress+":"+d.TargetPort, time.Duration(timeoutSecond*1e6))
		//conn, err = net.Dial("tcp", d.TargetAddress + ":" + d.TargetPort)
	} else {
		conn, err = net.Dial("tcp", d.TargetAddress+":"+d.TargetPort)
	}
	repeatTime := 1
	if d.IsRepeat {
		repeatTime = d.RepeatTime
	}

	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			cdata.timeoutCount += repeatTime
		} else {
			cdata.failCount += repeatTime
		}
		log.Error("StartSingleThread Dial err:" + err.Error())
		return
	}
	defer conn.Close()
	*outb = true

	waitResponseCh := make(chan int)

	tempCDataRead := getDefaultTotalCData()
	startSingleRead(conn, d, &tempCDataRead, waitResponseCh)

	for i := 0; i < repeatTime; i++ { // 重复发送
		var tempCData = getDefaultTotalCData()
		var ret bool
		startSingleSend(conn, d, &tempCData, &ret)
		mergeSingleResult(tempCData, cdata, ret)
		if d.IntervalTime > 0 && i != repeatTime-1 { // 最后一次不用休眠
			time.Sleep(time.Duration(d.IntervalTime * 1e6))
		}
	}
	if d.HasResponse {
		<-waitResponseCh
	}
	// 响应值取平均值，因为无法区分单个请求的响应
	// 减去因为发送间隔时间而阻塞的时间
	//tempCDataRead.totalResponseTime = (tempCDataRead.totalResponseTime - float64(d.IntervalTime * (repeatTime - 1))) / float64(d.ThreadNum) / float64(repeatTime)
	mergeSingleResult(tempCDataRead, cdata, true)
	if tempCDataRead.isTimeout { // 接收程序中，只要出现了timeout，则全判timeout(不覆盖 startSingleSend 中的 failCount )
		n := repeatTime - cdata.failCount
		cdata.timeoutCount += n
		cdata.succCount -= n
	} else if tempCDataRead.isFail { // 此处可能是服务端关闭连接，由于我方读取等待时间太长，而对端写完了数据就马上关闭了

	}
	/*if (d.HasResponse) {
		cdata.totalCostTime += cdata.totalResponseTime
	}*/
}

func startSingleRead(conn net.Conn, d model.CreateTaskData, outdata *countData, waitResponseCh chan<- int) {
	if d.HasResponse { // 读取数据
		go func() {
			var responseStart time.Time
			//var responseEnd time.Time
			defer close(waitResponseCh)
			count := 0
			var buf = make([]byte, default_read_buf_size)
			var preResponseTime time.Time
			repeatTime := 1
			if d.IsRepeat {
				repeatTime = d.RepeatTime
			}
			totalExpectBytes := d.ExpectedBytes * repeatTime
			cdata := outdata
			responseStart = time.Now()
			conn.SetReadDeadline(responseStart.Add(time.Duration(d.ReadTimeout * 1e6)))
			for { // 不计算最后一次超时时间
				preResponseTime = time.Now()
				responseStart = time.Now()
				if count > 0 { // 大于0时，要算上发送间隔时间
					conn.SetReadDeadline(preResponseTime.Add(time.Duration(d.ReadTimeout*1e6 + d.IntervalTime*1e6)))
				}
				rdn, err := conn.Read(buf) // 有发送间隔，此处至少会阻塞间隔那么长时间
				cdata.reciveDataSize += rdn
				//fmt.Println("recv:", string(buf))
				//responseEnd = time.Now()
				buf = make([]byte, default_read_buf_size)
				if err != nil {
					if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
						// 根据读超时时间自动断开的连接，不算错误，不计算本次响应时间
						if count > 0 {
							//responseEnd = preResponseTime
							cdata.isWaitReadEnd = true
							//cdata.totalResponseTime += float64(util.GetTimestamp(responseEnd) - util.GetTimestamp(responseStart))
							break
						} else { // 如果一次数据都没接收，超时了，说明是真的超时，这时候报错，不计算其超时时间
							log.Error("StartSingleThread Read timeout err:" + err.Error())
							cdata.isTimeout = true // 一次超时，后面的请求全算超时
						}
						return
					}
					cdata.isFail = true // 一次读取非超时错误，后面的请求全算错误
					log.Error("StartSingleThread Read err:" + err.Error())
					return
				}
				//cdata.totalResponseTime += float64(util.GetTimestamp(responseEnd) - util.GetTimestamp(responseStart))
				if d.ExpectedBytes > 0 && cdata.reciveDataSize >= totalExpectBytes {
					break
				}
				count++
			}
		}()
	}
}

func startSingleSend(conn net.Conn, d model.CreateTaskData, outdata *countData, outb *bool) {
	cdata := outdata
	*outb = false
	var err error

	start := time.Now()
	var end time.Time
	defer func() {
		end = time.Now()
		cdata.totalCostTime = float64(util.GetTimestamp(end) - util.GetTimestamp(start))
	}()
	//conn.SetWriteDeadline(time.Now().Add(time.Duration(5 * 1e9)))
	wrn, err := conn.Write(d.Data)
	if err != nil {
		cdata.failCount++
		log.Error("StartSingleThread Write err:" + err.Error())
		return
	}
	cdata.sendDataSize += wrn
	cdata.succCount++
	*outb = true
}

func startTestAndCompute(c model.CreateTaskData, outTotalCData *countData, syncCh chan int, mutex *sync.Mutex) {
	ch := make(chan int)
	var tempCData = getDefaultTotalCData()
	var ret bool
	startTime := time.Now()
	go startSingleThread(c, &tempCData, &ret, ch)
	// 计算结果
	<-ch
	endTime := time.Now()
	tempCData.threadCostTime = util.GetTimestamp(endTime) - util.GetTimestamp(startTime)
	mutex.Lock()
	mergeThreadResult(tempCData, outTotalCData, ret)
	mutex.Unlock()
	syncCh <- 0
}

// 开始任务，异步执行
func StartTask(taskId string) error {
	v, ok := taskMap[taskId]
	if !ok {
		return errors.New("ID为[" + taskId + "]的任务不存在")
	}
	if v.State != model.NOT_START {
		return errors.New("该任务正在运行")
	}

	go func() {
		startTime := time.Now()
		v.State = model.RUNNING
		v.StartTime = util.GetDateToStr(startTime, util.TIME_TEMPLATE_1)

		totalCData := getDefaultTotalCData()
		var lock sync.Mutex
		ch := make(chan int)
		for i := 0; i < v.ThreadNum; i++ {
			go startTestAndCompute(v.CreateTaskData, &totalCData, ch, &lock)
		}
		for i := 0; i < v.ThreadNum; i++ {
			<-ch
		}
		endTime := time.Now()
		totalCData.totalCostTime = float64(util.GetTimestamp(endTime) - util.GetTimestamp(startTime))
		// 计算各类平均值
		v.TotalRealCostTime = totalCData.totalCostTime
		if totalCData.isWaitReadEnd { // 如果发生了正常等待则请求总时间要减去这个等待时间
			totalCData.totalCostTime -= float64(v.ReadTimeout)
		}
		v.EndTime = util.GetDateToStr(endTime, util.TIME_TEMPLATE_1)
		v.TotalRequestCount = totalCData.totalRequestCount
		v.RequestAverageCostTime = totalCData.totalCostTime / float64(totalCData.totalRequestCount)
		v.RequestCostMaxTime = totalCData.costMaxTime
		v.RequestCostMinTime = totalCData.costMinTime
		//v.RequestAverageResponseTime = totalCData.totalResponseTime / float64(totalCData.succCount)
		//v.RequestResponseMaxTime = totalCData.responseMaxTime
		//v.RequestResponseMinTime = totalCData.responseMinTime
		if totalCData.succCount > 0 {
			v.ThreadAverageCostTime = float64(totalCData.threadCostTime) / float64(totalCData.succCount)
		}
		v.ThreadCostMaxTime = totalCData.threadCostMaxTime
		v.ThreadCostMinTime = totalCData.threadCostMinTIme
		if totalCData.totalCostTime == 0 {
			v.TransactionRate = 0
			v.Throughput = 0
		} else {
			v.TransactionRate = float64(totalCData.totalRequestCount) / (totalCData.totalCostTime / 1000)
			v.Throughput = float64(totalCData.sendDataSize+totalCData.reciveDataSize) / (totalCData.totalCostTime / 1000)
		}
		v.SuccTransactions = totalCData.succCount
		v.FailTransactions = totalCData.failCount
		v.TimeOutTransactions = totalCData.timeoutCount
		v.DataTransferred = totalCData.reciveDataSize + totalCData.sendDataSize
		v.SendBytes = totalCData.sendDataSize
		v.RecvBytes = totalCData.reciveDataSize
		v.TotalCostTime = totalCData.totalCostTime
		v.State = model.FINISH
	}()
	return nil
}
