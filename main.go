package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github/riny/go-grpc/user-system/config"
	"github/riny/go-grpc/user-system/repository"
	"github/riny/go-grpc/user-system/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	var grpcServerEndpoint = "localhost:50051"

	// 创建数据库连接
	var db, err = repository.NewDatabase(config.Dsn)
	if err != nil {
		log.Fatalf("Failed to connection database: %v", err)
	}
	defer func(db *repository.Database) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	// 创建user database connection
	var ur = repository.NewUserRepository(db)

	// 启动gRPC服务器
	go func() {
		if err := rungRPCServer(grpcServerEndpoint, ur); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	// 启动gRPC gateway服务器
	if err := runGatewayServer(grpcServerEndpoint, ur); err != nil {
		log.Fatalf("Failed to run gRPC-Gateway server: %v", err)
	}
}

func rungRPCServer(endpoint string, ur *repository.UserRepo) error {
	// 创建gRPC服务
	var grpcServer = grpc.NewServer()

	// 注册user-system gRPC服务
	service.RegisterUserServiceServer(grpcServer, service.NewImplementedUserServiceServer(ur))

	// 启动gRPC服务器
	var listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}
	return grpcServer.Serve(listener)
}

func runGatewayServer(grpcServerEndpoint string, ur *repository.UserRepo) error {
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
	return http.ListenAndServe("localhost:8080", allowAuthorizationMiddleware(mux, ur))
}

func allowAuthorizationMiddleware(h http.Handler, ur *repository.UserRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/api/v1/u/get" {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if exists, err := ur.CheckValidToken(nil, auth[7:]); !exists {
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
