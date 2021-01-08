package tcp

import (
	"fmt"
	"net"
	"strconv"
	"stress-tool/model"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	data := model.CreateTaskData{
		TargetAddress: "localhost",
		TargetPort:    "8080",
		ReadTimeout:   5000,
		ThreadNum:     1,
		IsRepeat:      false,
		RepeatTime:    0,
		IntervalTime:  0,
		HasResponse:   true,
		DataTypeMap: []map[string]int{
			{"0~3": model.STRING},
			{"4~5": model.STRING},
			{"6~7": model.STRING},
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0},
	}
	_, err := CreateTask(data)
	t.Logf("TestCreateTask expected: [%s], get: [%v]", "nil", err)
	if err != nil {
		t.Errorf("Test CreateTask fail")
	}

	data = model.CreateTaskData{
		TargetAddress: "localhost",
		TargetPort:    "8080",
		ThreadNum:     1,
		IsRepeat:      false,
		RepeatTime:    0,
		IntervalTime:  0,
		HasResponse:   true,
		DataTypeMap: []map[string]int{
			{"0~3": model.STRING},
			{"4~5": model.STRING},
			{"6~7": model.STRING},
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, // 数据多一位
	}
	_, err = CreateTask(data)
	t.Logf("TestCreateTask expected: [%s], get: [%v]", "not nil", err)
	if err == nil {
		t.Errorf("Test CreateTask fail")
	}
}

func TestStartSingleThread(t *testing.T) {
	serverResponse := "123456789."
	clientSend := "123456789."
	/*for i := 0; i < 12; i++ {
		serverResponse += serverResponse
		clientSend += clientSend
	}*/
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
		TargetAddress: "localhost",
		TargetPort:    strconv.Itoa(port),
		Timeout:       10000,
		ReadTimeout:   2000,
		ExpectedBytes: len(serverResponse),
		ThreadNum:     1,
		IsRepeat:      true,
		RepeatTime:    4,
		IntervalTime:  0,
		HasResponse:   true,
		DataTypeMap:   nil,
		Data:          []byte(clientSend),
	}

	var outCdata countData
	var outBool bool
	fmt.Println("start StartSingleThread...")
	go startSingleThread(data, &outCdata, &outBool, ch)
	<-ch
	<-ch2
	<-ch1
	fmt.Printf("%+v \n", outCdata)
	if !outBool {
		t.Error("Test StartSingleThread fail")
	}
}

func TestStartTask(t *testing.T) {
	serverResponse := "123456789."
	clientSend := "123456789.123456789."
	isRepeat := true
	repeatCount := 1
	threadNum := 1
	intervalTime := 0
	/*for i := 0; i < 10; i++ {
		serverResponse += serverResponse
		clientSend += clientSend
	}*/
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
		count = threadNum
	}
	go func() {
		for i := 0; i < threadNum; i++ {
			conn, err := listener.Accept()
			if err != nil {
				t.Error(err)
				return
			}
			go func() {
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

				for i := 0; i < repeatCount; i++ {
					_, err = conn.Write([]byte(serverResponse))
					// fmt.Println("test write ", n, " byte")
					if err != nil {
						t.Error(err)
					}
					time.Sleep(time.Duration(intervalTime * 1e6))
				}
				ch1 <- 0
			}()
		}
	}()

	fmt.Println("wait accept...")
	time.Sleep(300 * 1e6)
	// 构造测试数据
	data := model.CreateTaskData{
		TargetAddress: "localhost",
		TargetPort:    strconv.Itoa(port),
		Timeout:       2000,
		ReadTimeout:   500,
		ExpectedBytes: len(serverResponse),
		ThreadNum:     threadNum,
		IsRepeat:      isRepeat,
		RepeatTime:    repeatCount,
		IntervalTime:  intervalTime,
		HasResponse:   true,
		DataTypeMap: []map[string]int{
			{"0~999999": 2},
		},
		Data: []byte(clientSend),
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
