package code.barryyan.java_client;

import java.util.concurrent.TimeUnit;

import io.grpc.Channel;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.StatusRuntimeException;
import code.barryyan.proto.HelloRequest;
import code.barryyan.proto.HelloResponse;
import code.barryyan.proto.HelloServiceGrpc;

/**
 * @author  : YanMingxin
 * @email   : 1712229564@qq.com
 * @time    : 2023/3/4 19:20:12
 */ 
 
public class ClientClazz {

    private final HelloServiceGrpc.HelloServiceBlockingStub blockingStub;

    public ClientClazz(Channel channel) {
        blockingStub = HelloServiceGrpc.newBlockingStub(channel);
    }

    public static void main(String[] args) throws Exception {
        String target = "localhost:6666";

        ManagedChannel channel = ManagedChannelBuilder.forTarget(target)
                .usePlaintext()
                .build();

        HelloRequest request = HelloRequest
                .newBuilder()
                .setName("Barry Yan")
                .build();

        ClientClazz client = new ClientClazz(channel);
        HelloResponse response = client.blockingStub.sayHello(request);
        System.out.println(response.getMessage());

        channel.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
    }
}
