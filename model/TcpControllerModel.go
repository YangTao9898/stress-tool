package model

import "net"

type TcpConnRequest struct {
	TargetAddress string `json:"targetAddress"` // ip 或域名
	TargetPort    string `json:"targetPort"`
}

type InputDataMap struct {
	Type     string `json:"type"`     // 数据类型
	Length   string `json:"length"`   // 数据字节数
	Data     string `json:"data"`     // 数据
	IsBigEnd bool   `json:"isBigEnd"` // 数据大小端标志, true为大端，false为小端
}

type CreateTaskRequest struct {
	TargetAddress string         `json:"targetAddress"` // ip 或域名
	TargetPort    string         `json:"targetPort"`
	Timeout       string         `json:"timeout"`       // 超时时间，单位s, 小于等于 0 则设置默认超时时间
	ReadTimeout   string         `json:"readTimeout"`   // 读取超时时间，单位秒，tcp不会自动结束线程，需要设置一个读取超时时间
	ExpectedBytes string         `json:"expectedBytes"` // 读取到此字节数的数据则关闭连接，为 0 或负数则忽略此字段
	ThreadNum     string         `json:"threadNum"`
	IsRepeat      bool           `json:"isRepeat"`
	RepeatTime    string         `json:"repeatTime"`   // 重复发送次数，总次数等于 ThreadNum * RepeatTime
	IntervalTime  string         `json:"intervalTime"` // 重复发送间隔时间，单位毫秒，为0则不间断发送
	HasResponse   bool           `json:"hasResponse"`  // 是否有返回值
	DataMapArr    []InputDataMap `json:"dataMapArr"`
	SaveTaskTag   string         `json:"saveTaskTag"` // 存储任务名称
}

type GetAllTaskDescRequest struct {
	State int `json:"state"`
}

type TaskIdParamRequest struct {
	TaskId string `json:"taskId"`
}

type GetTaskDetailDataMap struct {
	Type       string   `json:"type"`       // 数据类型
	Length     string   `json:"length"`     // 数据字节数
	Data       string   `json:"data"`       // 数据
	IsBigEnd   bool     `json:"isBigEnd"`   // 数据大小端标志, true为大端，false为小端
	BinaryData string   `json:"binaryData"` // 对应的二进制数据
	ByteData   []string `json:"byteData"`   // 对应的八进制
}
type GetTaskDetailResponse struct {
	TargetAddress          string                 `json:"targetAddress"` // ip 或域名
	TargetPort             string                 `json:"targetPort"`
	Timeout                string                 `json:"timeout"`       // 超时时间，单位ms, 小于等于 0 则设置默认超时时间
	ReadTimeout            string                 `json:"readTimeout"`   // 读取超时时间，单位ms，tcp不会自动结束线程，需要设置一个读取超时时间
	ExpectedBytes          string                 `json:"expectedBytes"` // 读取到此字节数的数据则关闭连接，为 0 或负数则忽略此字段
	ThreadNum              string                 `json:"threadNum"`
	IsRepeat               bool                   `json:"isRepeat"`
	RepeatTime             string                 `json:"repeatTime"`   // 重复发送次数，总次数等于 ThreadNum * RepeatTime
	IntervalTime           string                 `json:"intervalTime"` // 重复发送间隔时间，单位毫秒，为0则不间断发送
	HasResponse            bool                   `json:"hasResponse"`  // 是否有返回值
	DataMapArr             []GetTaskDetailDataMap `json:"dataMapArr"`
	Taskid                 string                 `json:"taskid"` // 任务标识
	State                  string                 `json:"state"`  // 运行状态，参考上面的枚举值 NOT_START = 0 READY = 1 RUNNING = 2 FINISH = 3
	StartTime              string                 `json:"startTime"`
	EndTime                string                 `json:"endTime"`
	TotalRequestCount      string                 `json:"totalRequestCount"`      // 请求总次数
	RequestAverageCostTime string                 `json:"requestAverageCostTime"` // 单次请求花费的平均时间，单位 ms
	RequestCostMaxTime     string                 `json:"requestCostMaxTime"`     // 单次请求花费的最大时间
	RequestCostMinTime     string                 `json:"requestCostMinTime"`     // 单次请求花费的最小时间
	//RequestAverageResponseTime string `json:"requestAverageResponseTime"` // 单次请求响应的平均时间
	//RequestResponseMaxTime     string `json:"requestResponseMaxTime"` // 单次请求最大响应时间
	//RequestResponseMinTime     string `json:"requestResponseMinTime"` // 单次请求最小响应时间
	ThreadAverageCostTime string `json:"threadAverageCostTime"` // 线程执行完花费的平均时间
	ThreadCostMaxTime     string `json:"threadCostMaxTime"`     // 单个线程执行完花费的最大时间
	ThreadCostMinTime     string `json:"threadCostMinTime"`     // 单个线程执行完花费的最小时间
	TransactionRate       string `json:"transactionRate"`       // 每秒处理请求
	SuccTransactions      string `json:"succTransactions"`      // 请求成功数
	FailTransactions      string `json:"failTransactions"`      // 请求失败数
	TimeOutTransactions   string `json:"timeOutTransactions"`   // 请求超时数
	DataTransferred       string `json:"dataTransferred"`       // 总传输数据量
	Throughput            string `json:"throughput"`            // 每秒钟传输的数据量，吞吐量
	RecvBytes             string `json:"recvBytes"`             // 接收的总数据量
	TotalCostTime         string `json:"totalCostTime"`         // 不算失败数，总花费时间 ms
	TotalRealCostTime     string `json:"totalRealCostTime"`     // 实际花费时间
}

type SaveTcpTaskFileItem struct {
	SaveTaskId string            `json:"saveTaskId"`
	SaveTime   string            `json:"saveTime"`
	TaskData   CreateTaskRequest `json:"taskData"`
}

type SaveTcpTaskFileDesc struct {
	SaveTaskId    string `json:"saveTaskId"`
	SaveTaskTag   string `json:"saveTaskTag"`
	SaveTime      string `json:"saveTime"`
	TargetAddress string `json:"targetAddress"`
	TargetPort    string `json:"targetPort"`
	ThreadNum     string `json:"threadNum"`
}

type SaveTaskIdStruct struct {
	SaveTaskId string `json:"saveTaskId"`
}

type SaveTaskIdArrStruct struct {
	SaveTaskIdArr []string `json:"saveTaskIdArr"`
}

type TcpReturnTestConnStruct struct {
	Conn           net.Conn
	IsConn         bool
	DataMapArrList [][]InputDataMap `json:"dataMapArrList"` // 存储请求队列
}

type TcpTestReturnConnectResponse struct {
	Msg    string `json:"msg"`
	Result bool   `json:"result"`
}

type NumberRequest struct {
	Num int `json:"num"`
}

type RequestQueueUpdateDataRequest struct {
	QueueIndex int            `json:"queueIndex"`
	DataMapArr []InputDataMap `json:"dataMapArr"`
}
