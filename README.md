## 《玩转gRPC框架》专栏相关demo

有问题可以留言哈~

### 1 专栏信息

作者：https://blog.csdn.net/Mr_YanMingXin

E-mail：1712229564@qq.com

地址：https://blog.csdn.net/mr_yanmingxin/category_12172887.html

### 2 具体Demo列表

- [grpc-etcd](./grpc-etcd)：引入etcd服务注册中心，需要运行etcd
- [go-grpc-demo](./go-grpc-demo)：Go实现gRPC服务，可以调用Java提供的gRPC服务和为Java客户端提供服务
- [java-grpc-demo](./java-grpc-demo)：Java实现gRPC服务，可以调用Go提供的gRPC服务和为Go客户端提供服务
- [go-socket](./go-socket)：Go实现Socket
- [go-rpc](./go-rpc)：Go实现RPC调用
- [barry_rpc](./barry_rpc)：动手实现自己的RPC框架（TODO）
- [barry_lb](./barry_lb)：自定义负载均衡器

### 3 运行环境及方式

#### 3.1 涉及环境要求

JDK1.8+

Go1.17+

Maven 3.x

etcd 3.x

#### 3.2 具体运行方式

Java Demo

```shell
mvn clean 

mvn compile
```

Go Demo

```shell
go mod tidy 

go run xxx.go
```

### 4 注意点
