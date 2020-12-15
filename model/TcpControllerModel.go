package model

type TestConnectivityRequest struct {
	TargetAddress string `json:"targetAddress"` // ip 或域名
	TargetPort    string `json:"targetPort"`
}

type CreateTaskRequest struct {
	TargetAddress string         `json:"targetAddress"` // ip 或域名
	TargetPort    string         `json:"targetPort"`
	Timeout       string         `json:"timeout"`       // 超时时间，单位s, 小于等于 0 则设置默认超时时间
	ReadTimeout   string         `json:"readTimeout"`   // 读取超时时间，单位秒，tcp不会自动结束线程，需要设置一个读取超时时间
	ExpectedBytes string         `json:"expectedBytes"` // 读取到此字节数的数据则关闭连接，为 0 或负数则忽略此字段
	ThreadNum     string         `json:"threadNum"`
	IsRepeat      bool           `json:"isRepeat"`
	RepeatTime    string         `json:"repeatTime"`    // 重复发送次数，总次数等于 ThreadNum * RepeatTime
	IntervalTime  string         `json:"intervalTime"`  // 重复发送间隔时间，单位毫秒，为0则不间断发送
	HasResponse   bool           `json:"hasResponse"`   // 是否有返回值
	DataMapArr       []struct {
		Type   string `json:"type"` // 数据类型
		Length string    `json:"length"` // 数据字节数
		Data   string `json:"data"` // 数据
	} `json:"dataMapArr"`
}

type GetAllTaskDescRequest struct {
	State int `json:"state"`
}