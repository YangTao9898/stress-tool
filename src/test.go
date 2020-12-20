package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func listen() {
	listener, e := net.Listen("tcp", "localhost:14000")
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println("listen on 14000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		defer conn.Close()
		fmt.Println("Accept conn on:" + conn.RemoteAddr().String())

		go func() {
			n, err := conn.Write([]byte("Hello, Got your."))
			fmt.Println("test write ", n, " byte")
			if err != nil {
				fmt.Println(err)
			}

			bs := make([]byte, 1024)
			for {
				n, err := conn.Read(bs)
				fmt.Println("recive ", n, " byte:", string(bs))
				bs = make([]byte, 1024)
				if err != nil {
					fmt.Println("Dissconnected on:" + conn.RemoteAddr().String())
					return
				}
			}
		}()
	}
}

func main() {
	//listen()
	/*for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				util.GetLogger().Debug(strconv.Itoa(i) + ":   123456789 123456789 123456789")
			}
		}()
	}
	select {
	}*/

	/*dataInt, err := strconv.ParseFloat("256", 10)
	fmt.Println(err)
	fmt.Println(dataInt)*/
	/*s1 := []string{"123"}
	s2 := []string{"334", "556"}
	s1 = append(s1, s2...)
	fmt.Println(s1)*/
	var a int64 = 1
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, a)
	fmt.Println(buf.Len(), buf.Bytes())
}
