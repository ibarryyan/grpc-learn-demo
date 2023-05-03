package t_test

import (
	"fmt"
	"net"
	"testing"
)

func TestSocketServer(t *testing.T) {
	// 监听端口
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started, listening on port 8888")

	// 循环等待客户端连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// 启动一个 goroutine 处理客户端请求
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// 循环读取客户端发送的数据
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}
		// 处理客户端发送的数据
		fmt.Println("Received from", conn.RemoteAddr(), ":", string(buf[:n]))
		// 回复客户端
		conn.Write([]byte("Message received."))
	}
}


func TestSocketClient(t *testing.T) {
	// 连接服务端
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	// 发送给服务端
	conn.Write([]byte("hello "))

	// 接收服务端回复并打印
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Received from server:", string(buf[:n]))
}