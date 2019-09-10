package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := "9090"
	start(port)
}

//启动服务器
func start(port string) {
	host := ":" + port
	//获取tcp地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Printf("Resolve tcp addr failed:%v\n", err)
		return
	}

	//监听
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("Listener tcp port failed:%v\n", err)
		return
	}

	//建立连接池，用于广播消息
	conns := make(map[string]net.Conn)
	//消息通道
	messageChan := make(chan string, 10)
	//广播消息
	go BoroadMessages(&conns, messageChan)

	//启动
	for {
		fmt.Printf("listener port %s ...\n", port)
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Accept failed:%v\n", err)
			continue
		}
		conns[conn.RemoteAddr().String()] = conn
		fmt.Println(conns)

		//处理消息
		go Handler(conn, &conns, messageChan)
	}
}

func Handler(conn net.Conn, conns *map[string]net.Conn, message chan string) {
	buf := make([]byte, 1024)
	for {
		lenth, err := conn.Read(buf)
		if err != nil {
			log.Printf("Read client message failed:%v\n", err)
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}
		recvStr := string(buf[0:lenth])
		message <- recvStr
	}
}

func BoroadMessages(conns *map[string]net.Conn, message chan string) {
	for {
		//不断从通道里读数据
		msg := <-message
		fmt.Println(msg)
		//向所有客户端发送信息
		for key, conn := range *conns {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("Broad message to %s failed:%v\n", key, err)
				delete(*conns, key)
			}
		}
	}
}
