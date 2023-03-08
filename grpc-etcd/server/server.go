package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"grpc-etcd/proto"

	"github.com/ozonru/etcd/clientv3"
	"google.golang.org/grpc"
)

var (
	cli         *clientv3.Client
	Schema      = "ns"
	Host        = "127.0.0.1"
	Port        = 3000             //端口
	ServiceName = "helloService"   //服务名称
	EtcdAddr    = "127.0.0.1:2379" //etcd地址
)

type HelloServer struct{}

func (s *HelloServer) Hello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello %s , I am Barry Yan", req.GetName()),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", Port))
	if err != nil {
		fmt.Println("Listen network err :", err)
		return
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Printf("err %s", err)
		}
	}()

	srv := grpc.NewServer()
	proto.RegisterHelloServiceServer(srv, &HelloServer{})

	err = register(EtcdAddr, ServiceName, fmt.Sprintf("%s:%d", Host, Port), 10)
	if err != nil {
		fmt.Printf("register err %s", err)
	}

	//关闭信号处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		unRegister(ServiceName, fmt.Sprintf("%s:%d", Host, Port))
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	if srv.Serve(listener) != nil {
		fmt.Println("rpc server err : ", err)
	}
}

func register(etcdAddr, serviceName, serverAddr string, ttl int64) error {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 50 * time.Second,
	})
	if err != nil {
		fmt.Printf("connection server err : %s\n", err)
		return err
	}
	//与etcd建立长连接，开启心跳检测
	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		for {
			resp, err := cli.Get(context.TODO(), getKey(serviceName, serverAddr))
			if err != nil {
				fmt.Printf("get server address err : %s", err)
			} else if resp.Count == 0 { //尚未注册
				if keepAlive(serviceName, serverAddr, ttl) != nil {
					fmt.Printf("keepAlive err : %s", err)
				}
			}
			<-ticker.C
		}
	}()
	return nil
}

func keepAlive(serviceName, serverAddr string, ttl int64) error {
	//创建租约
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		fmt.Printf("create grant err : %s\n", err)
		return err
	}
	//将服务地址注册到etcd中
	_, err = cli.Put(context.Background(), getKey(serviceName, serverAddr), serverAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		fmt.Printf("register service err : %s", err)
		return err
	}
	
	ch, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		fmt.Printf("KeepAlive err : %s\n", err)
		return err
	}
	
	go func() {
		for {
			<-ch
		}
	}()
	return nil
}

func unRegister(serviceName, serverAddr string) {
	if cli != nil {
		response, err := cli.Delete(context.Background(), getKey(serviceName, serverAddr))
		if err != nil {
			fmt.Println("Listen network err :", err, response)
		}
	}
}

func getKey(serviceName, serverAddr string) string {
	return "/" + Schema + "/" + serviceName + "/" + serverAddr
}
