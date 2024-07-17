package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

type Data struct {
	Msg string
}

func (t *Data) Byte() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(t)
	if err != nil {
		log.Fatal("encode error", err)
	}
	return buffer.Bytes()
}

func DataDecoder(data []byte) *Data {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	tra := &Data{}
	err := decoder.Decode(&tra)
	if err != nil {
		log.Fatal("decode error", err)
	}
	return tra
}

func ral(data *Data, tcpAddr *net.TCPAddr) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(data.Byte())
	if err != nil {
		panic(err)
	}
}

func server() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8204")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	log.Println("服务创建成功")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			// time.Sleep(1 * time.Second)
			// data := &Data{
			// 	Msg: "333",
			// }
			// _, err := conn.Write(data.Byte())
			// if err != nil {
			// 	fmt.Println("失败")
			// 	return
			// }
			// fmt.Println("成功")
			// return

			resv := make([]byte, 1024)
			n, err := conn.Read(resv)
			if err != nil {
				panic(err)
			}
			if n > 0 && n < 1025 {
				data := DataDecoder(resv)
				fmt.Println(data.Msg)
			} else {
				return
			}
		}(conn)
	}
}

func client() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "192.168.1.121:8204")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	conn.Close()
	data4 := &Data{
		Msg: "444",
	}
	_, err = conn.Write(data4.Byte())
	if err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Second)

	data1 := &Data{
		Msg: "111",
	}
	ral(data1, tcpAddr)
	fmt.Println("客户端第一个请求成功")

	time.Sleep(1 * time.Second)

	data2 := &Data{
		Msg: "222",
	}
	ral(data2, tcpAddr)
	fmt.Println("客户端第二个请求成功")
}

func main() {
	go server()
	time.Sleep(3 * time.Second)
	go client()
	time.Sleep(100 * time.Second)
}
