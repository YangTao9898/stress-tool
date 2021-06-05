package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s 9090(listen port) to run", os.Args[0])
		return
	}
	listener, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen error, cause: %s", err.Error())
		return
	}
	fmt.Println("listen succ on " + os.Args[1])

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept error, cause: %s", err.Error())
			return
		}
		fmt.Println("new conn accept")
		go func() {
			defer conn.Close()
			for {
				// 读
				buf := make([]byte, 1024)
				// 写
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "recv err, cause: %s", err.Error())
					return
				}
				readBuf := buf[:n]
				fmt.Println("recv and write: " + string(readBuf))
				_, err = conn.Write(readBuf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "write err, cause: %s", err.Error())
				}
			}
		}()
	}
}
