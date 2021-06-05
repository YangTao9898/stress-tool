# tcp压测工具

#### 介绍
本项目为tcp压测工具，可以对tcp服务端进行压力测试，同时也可以像postman一样对tcp服务端一样进行tcp的接口测试，发送指定的值获取结果，并且支持对结果进行转换。

#### 安装使用教程

#####需要安装go环境，如果不想自己打包，也可以直接下载打包完后的包
#####1. windows
`cd stress-tool `

`go build -o stress_tool.exe main.go`

`./stress_tool.exe 8080`
#####2. linux or macos
`cd stress-tool `

`go build -o stress_tool main.go`

`./stress_tool 8080`

###界面访问
浏览器输入 ip:8080 即可访问压测界面

#### 使用说明

1.  创建压测界面用于创建压测任务
2.  查看压测任务用于查看压测结果
3.  TCP返回测试界面用于测试TCP接口

####服务端模拟程序说明
本项目附带一个tcp服务端的模拟程序，接收到什么输入，就会返回什么输出，安装使用教程如下

`cd stress-tool/test`

`go build -o mockTcpServer mockTcpServer.go`

`./mockTcpServer 8888`
