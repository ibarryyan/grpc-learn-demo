/**
 * @author  : YanMingxin
 * @email   : 1712229564@qq.com
 * @time    : 2023/3/4 19:20:12
 */ 

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"grpc-etcd/proto"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ozonru/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	cli         *clientv3.Client
	Schema      = "ns"
	ServiceName = "helloService"   //服务名称
	EtcdAddr    = "127.0.0.1:2379" //etcd地址
)

func main() {
	r := NewEtcdResolver(EtcdAddr)
	resolver.Register(r)

	conn, err := grpc.Dial(r.Scheme()+"://barry/"+ServiceName, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("connect err : %s", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("close err : %s", err)
		}
	}()

	c := proto.NewHelloServiceClient(conn)

	resp, err := c.Hello(context.Background(), &proto.HelloRequest{Name: "gRPC"})
	if err != nil {
		fmt.Printf("call service err : %s", err)
		return
	}
	fmt.Printf("resp : %v", resp)
}

type EtcdResolver struct {
	etcdAddr   string
	clientConn resolver.ClientConn
}

func NewEtcdResolver(etcdAddr string) resolver.Builder {
	return &EtcdResolver{etcdAddr: etcdAddr}
}

func (r *EtcdResolver) Scheme() string {
	return Schema
}

func (r *EtcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	fmt.Println(rn)
}

func (r *EtcdResolver) Close() {}

//Build 构建解析器 grpc.Dial()时调用
func (r *EtcdResolver) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	//构建etcd client
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(r.etcdAddr, ";"),
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			fmt.Printf("connect etcd err : %s\n", err)
			return nil, err
		}
	}
	r.clientConn = clientConn
	go r.watch("/" + target.Scheme + "/" + target.Endpoint + "/")
	return r, nil
}

//watch机制：监听etcd中某个key前缀的服务地址列表的变化
func (r *EtcdResolver) watch(keyPrefix string) {
	var addrList []resolver.Address
	resp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("get service list err : ", err)
	} else {
		for i := range resp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(resp.Kvs[i].Key), keyPrefix)})
		}
	}
	r.clientConn.UpdateState(resolver.State{Addresses: addrList})
	//监听服务地址列表的变化
	rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exists(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.clientConn.UpdateState(resolver.State{Addresses: addrList})
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					r.clientConn.UpdateState(resolver.State{Addresses: s})
				}
			}
		}
	}
}

func exists(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
