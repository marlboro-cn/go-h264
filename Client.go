package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	//接收终端输入的数据
	Start(os.Args[1])
}

func Start(tcpAddrStr string) {
	//1.根据输入的IP加端口生成数据
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddrStr)
	if err != nil {
		log.Printf("Resolve tcp addr failed:%v\n", err)
		return
	}
	//2.建立链接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial to server failed:%v\n", err)
		return
	}
	//向服务器发送信息
	go SendMsg(conn)
	//从服务器接收信息
	buf := make([]byte, 1024)
	for {
		lenth, err := conn.Read(buf)
		if err != nil {
			log.Printf("Recv server msg failed:%v\n", err)
			conn.Close()
			os.Exit(0)
			break
		}
		fmt.Println(string(buf[0:lenth]))
	}
}

func SendMsg(conn net.Conn) {
	userName := "乔萝莉～"
	for {
		var input string
		fmt.Scanln(&input)
		if input == "/q" || input == "/quit" {
			fmt.Println("Byebye~")
			conn.Close()
			os.Exit(0)
		}
		if len(input) > 0 {
			msg := userName + ":" + input
			_, err := conn.Write([]byte(msg))
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}
