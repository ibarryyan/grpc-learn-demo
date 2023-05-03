package t_test

import (
	"fmt"
	"net"
	"net/rpc"
	"testing"
)

type HelloService struct {}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "Hello, " + request
	return nil
}

func TestRpcServer(t *testing.T) {
	// 创建RPC服务器
	rpcServer := rpc.NewServer()
	// 注册HelloService服务
	err := rpcServer.RegisterName("HelloService", new(HelloService))
	if err != nil {
		fmt.Println(err)
	}
	// 启动监听
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
	}
	// 接收连接并为每个连接创建一个goroutine处理RPC请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go rpcServer.ServeConn(conn)
	}
}

func TestRpcClient(t *testing.T) {
	// 连接RPC服务器
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()
	// 调用HelloService的Hello方法
	var reply string
	err = client.Call("HelloService.Hello", "world", &reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}
