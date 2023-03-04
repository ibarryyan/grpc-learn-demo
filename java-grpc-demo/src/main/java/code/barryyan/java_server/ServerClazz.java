package code.barryyan.java_server;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import code.barryyan.proto.HelloRequest;
import code.barryyan.proto.HelloResponse;
import code.barryyan.proto.HelloServiceGrpc;

/**
 * @author  : YanMingxin
 * @email   : 1712229564@qq.com
 * @time    : 2023/3/4 19:20:12
 */ 
 
public class ServerClazz {

    private final int SERVER_PORT = 50051;

    private Server server;

    private void start() throws IOException {
        server = ServerBuilder.forPort(SERVER_PORT)
                .addService(new HelloServiceGrpcImpl())
                .build()
                .start();
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            try {
                ServerClazz.this.stop();
            } catch (InterruptedException e) {
                e.printStackTrace(System.err);
            }
        }));
    }

    private void stop() throws InterruptedException {
        if (server != null) {
            server.shutdown().awaitTermination(30, TimeUnit.SECONDS);
        }
    }

    private void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        final ServerClazz server = new ServerClazz();
        server.start();
        server.blockUntilShutdown();
    }

    static class HelloServiceGrpcImpl extends HelloServiceGrpc.HelloServiceImplBase {
        @Override
        public void sayHello(HelloRequest req, StreamObserver<HelloResponse> responseObserver) {
            HelloResponse resp = HelloResponse.newBuilder()
                    .setMessage("Hello " + req.getName())
                    .build();
            System.out.println("server: Hello " + req.getName());
            responseObserver.onNext(resp);
            responseObserver.onCompleted();
        }
    }
}
