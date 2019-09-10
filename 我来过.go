package main

import (
	"fmt"
	"log"
	"net"
)

//聊天系统（服务器端）
func main() {
	port := "9090"
	Start(port)
}

//启动服务器
func Start(port string) {
	host := ":" + port
	//获取tcp地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Printf("resolve tcp addr failed:%v\n", err)
		return
	}
	//监听
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("listen tcp port failed:%v\n", err)
		return
	}
	//连接池
	conns := make(map[string]net.Conn)
	//消息通道
	messageChan := make(chan string, 10)
	// ⼴播消息
	go BroadMessages(&conns, messageChan)

	for {
		fmt.Printf("listen port %s ...\n", port)
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Accept failed:%v\n", err)
			continue
		}
		//扔进连接池
		conns[conn.RemoteAddr().String()] = conn
		fmt.Println(conns)
		// 处理消息
		go Handler(conn, &conns, messageChan)
	}
}

//广播
func BroadMessages(conns *map[string]net.Conn, messages chan string) {
	for {
		//不断从消息通道里读取消息
		msg := <-messages
		fmt.Println(msg)
		//广播消息
		for key, conn := range *conns {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("broad message to %s failed :%v\n", key, err)
				delete(*conns, key)
			}
		}
	}
}

//收消息
func Handler(conn net.Conn, conns *map[string]net.Conn, messages chan string) {
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("read client messages failed:%v\n", err)
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}
		//把收到的消息写入消息通道
		recvStr := string(buf[0:length])
		messages <- recvStr
	}
}
