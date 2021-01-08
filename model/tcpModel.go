package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const ( // 发送的数据类型
	NUMBER = 0
	FLOAT  = 1
	STRING = 2
)

type CreateTaskData struct {
	TargetAddress string `json:"targetAddress"` // ip 或域名
	TargetPort    string `json:"targetPort"`
	Timeout       int    `json:"timeout"`       // 超时时间，单位ms, 小于等于 0 则设置默认超时时间
	ReadTimeout   int    `json:"readTimeout"`   // 读取超时时间，单位ms，tcp不会自动结束线程，需要设置一个读取超时时间
	ExpectedBytes int    `json:"expectedBytes"` // 读取到此字节数的数据则关闭连接，为 0 或负数则忽略此字段
	ThreadNum     int    `json:"threadNum"`
	IsRepeat      bool   `json:"isRepeat"`
	RepeatTime    int    `json:"repeatTime"`   // 重复发送次数，总次数等于 ThreadNum * RepeatTime
	IntervalTime  int    `json:"intervalTime"` // 重复发送间隔时间，单位毫秒，为0则不间断发送
	HasResponse   bool   `json:"hasResponse"`  // 是否有返回值
	/*	NUMBER     = 0 FLOAT   = 1 STRING  = 2
	*	数据类型 key 为数据索引范围，value 为数据类型，如上枚举所叙述，key:0~3 value:0 key:4~15 value:2 则表示前四个字节是NUMBER，后12个字节为string
	*   为了使之有序，每一个元素用一个map来装，形成一个字典数组
	 */
	DataTypeMap []map[string]int `json:"dataTypeMap"`
	IsBigEnd    []bool           `json:"isBigEnd"` // 数字类型数据的存储格式，true为大端 false为小端
	Data        []byte           `json:"data"`     // 请求数据
}

// 将形如 "0~3" 的字符串分割成数字 0 和 3
func CreateTaskDataKeySplite(v string) (int, int, error) {
	nums := strings.Split(v, "~")
	if len(nums) != 2 {
		return -1, -1, errors.New(fmt.Sprintf("数据范围须遵循 0~4 格式, 但是其格式为:[%s]", v))
	}
	var arr [2]int
	var err1, err2 error
	arr[0], err1 = strconv.Atoi(nums[0])
	arr[1], err2 = strconv.Atoi(nums[1])
	if err1 != nil || err2 != nil {
		return -1, -1, errors.New(fmt.Sprintf("expect number, but get [%s],[%s]", nums[0], nums[1]))
	}
	return arr[0], arr[1], nil
}

const ( // 任务状态
	NOT_START = 0
	READY     = 1
	RUNNING   = 2
	FINISH    = 3
)

type TaskDealData struct {
	CreateTaskData         `json:"createTaskData"`
	Taskid                 string  `json:"taskid"` // 任务标识
	State                  int     `json:"state"`  // 运行状态，参考上面的枚举值 NOT_START = 0 READY = 1 RUNNING = 2 FINISH = 3
	StartTime              string  `json:"startTime"`
	EndTime                string  `json:"endTime"`
	TotalRequestCount      int     `json:"totalRequestCount"`      // 请求总次数
	RequestAverageCostTime float64 `json:"requestAverageCostTime"` // 单次请求花费的平均时间，单位 ms
	RequestCostMaxTime     float64 `json:"requestCostMaxTime"`     // 单次请求花费的最大时间
	RequestCostMinTime     float64 `json:"requestCostMinTime"`     // 单次请求花费的最小时间
	//RequestAverageResponseTime float64 `json:"requestAverageResponseTime"` // 单次请求响应的平均时间
	//RequestResponseMaxTime     float64 `json:"requestResponseMaxTime"` // 单次请求最大响应时间
	//RequestResponseMinTime     float64 `json:"requestResponseMinTime"` // 单次请求最小响应时间
	ThreadAverageCostTime float64 `json:"threadAverageCostTime"` // 线程执行完花费的平均时间
	ThreadCostMaxTime     int64   `json:"threadCostMaxTime"`     // 单个线程执行完花费的最大时间
	ThreadCostMinTime     int64   `json:"threadCostMinTime"`     // 单个线程执行完花费的最小时间
	TransactionRate       float64 `json:"transactionRate"`       // 每秒处理请求
	SuccTransactions      int     `json:"succTransactions"`      // 请求成功数
	FailTransactions      int     `json:"failTransactions"`      // 请求失败数
	TimeOutTransactions   int     `json:"timeOutTransactions"`   // 请求超时数
	DataTransferred       int     `json:"dataTransferred"`       // 总传输数据量
	SendBytes             int     `json:"sendBytes"`
	RecvBytes             int     `json:"recvBytes"`
	Throughput            float64 `json:"throughput"`        // 每秒钟传输的数据量，吞吐量
	TotalCostTime         float64 `json:"totalCostTime"`     // 总花费时间 ms，不算最后一次等待的超时时间
	TotalRealCostTime     float64 `json:"totalRealCostTime"` // 实际花费时间
}
type TaskDealDataDescript struct {
	TaskId        string `json:"taskId"`
	TargetAddress string `json:"targetAddress"`
	TargetPort    string `json:"targetPort"`
	State         int    `json:"state"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
}
