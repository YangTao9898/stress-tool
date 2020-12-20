package tcp

import (
	"fmt"
	"net"
	"strconv"
	"stress-tool/model"
	"testing"
	"time"
)

var testValues [][]map[string]int
var testResult []bool

func addTestItem(dataTypeMap map[string]int, result bool) {
	arr := make([]map[string]int, len(dataTypeMap))
	index := 0
	for k, v := range dataTypeMap {
		m := make(map[string]int, 1)
		m[k] = v
		arr[index] = m
		index++
	}
	testValues = append(testValues, arr)
	testResult = append(testResult, result)
}

func TestCheckDataRangeIsContinuous(t *testing.T) {
	addTestItem(map[string]int{
		"0~3": model.STRING,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, true)
	addTestItem(map[string]int{
		"1~3": model.STRING,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"-1~3": model.STRING,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"3~0": model.STRING,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"3~0": model.STRING,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~3": model.STRING,
		"4~-12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~3": model.STRING,
		"4~12": model.STRING,
		"11~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~11": model.STRING,
		"12~12": model.STRING,
		"13~20": model.STRING,
	}, true)
	addTestItem(map[string]int{
		"0~4": model.NUMBER,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~4": model.FLOAT,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~4": model.FLOAT,
		"4~12": model.STRING,
		"13~20": model.STRING,
	}, false)

	i := 0
	var msg string
	for length := len(testValues); i < length; i++ {
		m := testValues[i]
		_, msg = CheckDataRangeIsContinuous(m)
		expected := ""
		if (!testResult[i]) {
			expected = "Not null string"
		}
		t.Logf("Test CheckDataRangeIsContinuous: param: [%v], expected: [%s], get: [%s]", m, expected, msg)
		if msg == "" && !testResult[i] || msg != "" && testResult[i] {
			t.Errorf("TestCheckDataRangeIsContinuous fail")
		}
	}
}

func TestCreateTask(t *testing.T) {
	data := model.CreateTaskData{
		TargetAddress: "localhost",
		TargetPort:   "8080",
		ReadTimeout:  5000,
		ThreadNum:    1,
		IsRepeat:     false,
		RepeatTime:   0,
		IntervalTime: 0,
		HasResponse:  true,
		DataTypeMap:  []map[string]int{
			{"0~3": model.STRING},
			{"4~5": model.STRING},
			{"6~7": model.STRING},
		},
		Data:         []byte{0, 0, 0, 0, 0, 0, 0, 0},
	}
	_, err := CreateTask(data)
	t.Logf("TestCreateTask expected: [%s], get: [%v]", "nil", err)
	if err != nil {
		t.Errorf("Test CreateTask fail")
	}

	data = model.CreateTaskData{
		TargetAddress: "localhost",
		TargetPort:   "8080",
		ThreadNum:    1,
		IsRepeat:     false,
		RepeatTime:   0,
		IntervalTime: 0,
		HasResponse:  true,
		DataTypeMap:  []map[string]int{
			{"0~3": model.STRING},
			{"4~5": model.STRING},
			{"6~7": model.STRING},
		},
		Data:         []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, // 数据多一位
	}
	_, err = CreateTask(data)
	t.Logf("TestCreateTask expected: [%s], get: [%v]", "not nil", err)
	if err == nil {
		t.Errorf("Test CreateTask fail")
	}
}

func TestStartSingleThread(t *testing.T) {
	serverResponse := "I am test server."
	clientSend := "Hello, I am test client."
	for i := 0; i < 12; i++ {
		serverResponse += serverResponse
		clientSend += clientSend
	}
	ch := make(chan int)
	ch1 := make(chan int)
	ch2 := make(chan int)

	// 启动测试专用服务
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Error(err)
		return
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Println("test server start on port:", port)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		fmt.Println("Accept conn")

		go func() {
			bs := make([]byte, 1024)
			for {
				_, err := conn.Read(bs)
				// fmt.Println("recive ", n, " byte:", string(bs))
				bs = make([]byte, 1024)
				if err != nil {
					ch2 <- 0
					return
				}
			}
		}()

		_, err = conn.Write([]byte(serverResponse))
		// fmt.Println("test write ", n, " byte")
		if err != nil {
			t.Error(err)
		}
		ch1 <- 0
	}()

	fmt.Println("wait accept...")
	time.Sleep(300 * 1e6)
	// 构造测试数据
	data := model.CreateTaskData{
		TargetAddress:   "localhost",
		TargetPort:     strconv.Itoa(port),
		Timeout:        10000,
		ReadTimeout:    1000,
		ExpectedBytes: len(serverResponse),
		ThreadNum:      5,
		IsRepeat:       false,
		RepeatTime:     0,
		IntervalTime:   0,
		HasResponse:    true,
		DataTypeMap:    nil,
		Data:           []byte(clientSend),
	}

	var outCdata countData
	var outBool bool
	fmt.Println("start StartSingleThread...")
	go StartSingleThread(data, &outCdata, &outBool, ch)
	<-ch
	<-ch2
	<-ch1
	fmt.Printf("%+v \n", outCdata)
	if (!outBool) {
		t.Error("Test StartSingleThread fail")
	}
}

func TestStartTask(t *testing.T) {
	serverResponse := "I am test server."
	clientSend := "Hello, I am test client."
	isRepeat := true
	repeatCount := 10
	for i := 0; i < 10; i++ {
		serverResponse += serverResponse
		clientSend += clientSend
	}
	//ch := make(chan int)
	ch1 := make(chan int)
	ch2 := make(chan int)

	// 启动测试专用服务
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Error(err)
		return
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Println("test server start on port:", port)

	count := 1
	if isRepeat {
		count = repeatCount
	}
	go func() {
		for i := 0; i < count; i++ {
			conn, err := listener.Accept()
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()
			//fmt.Println("Accept conn")

			go func() {
				bs := make([]byte, 1024)
				for {
					_, err := conn.Read(bs)
					// fmt.Println("recive ", n, " byte:", string(bs))
					bs = make([]byte, 1024)
					if err != nil {
						ch2 <- 0
						return
					}
				}
			}()

			_, err = conn.Write([]byte(serverResponse))
			// fmt.Println("test write ", n, " byte")
			if err != nil {
				t.Error(err)
			}
			ch1 <- 0
		}
	}()

	fmt.Println("wait accept...")
	time.Sleep(300 * 1e6)
	// 构造测试数据
	data := model.CreateTaskData{
		TargetAddress:   "localhost",
		TargetPort:     strconv.Itoa(port),
		Timeout:        10000,
		ReadTimeout:    1000,
		ExpectedBytes: len(serverResponse),
		ThreadNum:      5,
		IsRepeat:       isRepeat,
		RepeatTime:     repeatCount,
		IntervalTime:   0,
		HasResponse:    true,
		DataTypeMap:    []map[string]int{
			{"0~999999": 2},
		},
		Data:           []byte(clientSend),
	}

	taskId, err := CreateTask(data)
	if err != nil {
		t.Error(err)
	}
	err = StartTask(taskId)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("wait start...")
	for i := 0; i < count; i++ {
		<-ch2
		<-ch1
	}
	fmt.Println("wait end...")

	task := GetTaskByTaskId(taskId)
	task.CreateTaskData.Data = nil
	fmt.Printf("%+v \n", *task)

	if err != nil {
		t.Error(err)
	}
}

func TestConvertTaskDealDataToJson(t *testing.T) {
	cdata := model.CreateTaskData{
		TargetAddress:   "localhost",
		TargetPort:     "8080",
		Timeout:        10000,
		ReadTimeout:    1000,
		ExpectedBytes:  32,
		ThreadNum:      5,
		IsRepeat:       false,
		RepeatTime:     0,
		IntervalTime:   0,
		HasResponse:    true,
		DataTypeMap:    []map[string]int{
			{"0~999999": 2},
		},
		Data:           []byte("xxxxxxxxxx"),
	}
	td := model.TaskDealData{
		CreateTaskData:             cdata,
		Taskid:                     "20201126093420",
		State:                      3,
		StartTime:                  "2020-11-26 09:34:20",
		EndTime:                    "2020-11-26 09:34:21",
		TotalRequestCount:          100,
		RequestAverageCostTime:     0.12,
		RequestCostMaxTime:         2,
		RequestCostMinTime:         0,
		RequestAverageResponseTime: 0.01,
		RequestResponseMaxTime:     1,
		RequestResponseMinTime:     0,
		TransactionRate:            8333.333333333334,
		SuccTransactions:           100,
		FailTransactions:           0,
		TimeOutTransactions:        0,
		DataTransferred:            9900032,
		Throughput:                 825002666.6666666,
		TotalCostTime:              12,
	}

	s, e := ConvertTaskDealDataToJson(td)
	if e != nil {
		t.Error(e)
	}
	fmt.Println("ConvertTaskDealDataToJson result:", s)
}






















