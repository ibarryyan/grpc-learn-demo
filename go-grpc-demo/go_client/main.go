/*
 * @author  : YanMingxin
 * @email   : 1712229564@qq.com
 * @time    : 2023/3/4 19:20:12
 */ 

package main

import (
	"context"
	"log"
	
	"go-grpc-demo/proto/hello_proto"
	
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":6666", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("正在监听服务端 : %v\n", err)
	}
	defer conn.Close()
	client := hello_proto.NewHelloServiceClient(conn)

	req := &hello_proto.HelloRequest{
		Name: "Barry Yan",
	}

	resp, err := client.SayHello(context.TODO(), req)
	if err != nil {
		log.Fatalf("请求错误 : %v\n", err)
	}
	log.Printf("响应内容 : %v\n", resp)
}
