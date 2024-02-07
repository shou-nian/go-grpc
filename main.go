package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github/riny/go-grpc/user-system/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	var grpcServerEndpoint = "localhost:50051"

	// 启动gRPC服务器
	go func() {
		if err := rungRPCServer(grpcServerEndpoint); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	// 启动gRPC gateway服务器
	//go func() {
	//	if err := runGatewayServer(grpcServerEndpoint); err != nil {
	//		log.Fatalf("Failed to run gRPC-Gateway server: %v", err)
	//	}
	//}()
	if err := runGatewayServer(grpcServerEndpoint); err != nil {
		log.Fatalf("Failed to run gRPC-Gateway server: %v", err)
	}
}

func rungRPCServer(endpoint string) error {
	// 创建gRPC服务
	var grpcServer = grpc.NewServer()

	// 注册gRPC服务
	service.RegisterUserServiceServer(grpcServer, &service.ImplementedUserServiceServer{})

	// 启动gRPC服务器
	var listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}
	return grpcServer.Serve(listener)
}

func runGatewayServer(grpcServerEndpoint string) error {
	// 创建ctx和取消函数
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// 创建http分发对象
	var mux = runtime.NewServeMux()

	// 连接gRPC服务
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	err := service.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// 启动http服务
	return http.ListenAndServe("localhost:8080", mux)
}
