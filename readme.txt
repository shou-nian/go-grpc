需要在settings - protocol buffers中设置protoc路径：D:\tools\protoc-25.2-win64\include

//
protoc --grpc-gateway_out=. --go-grpc_out=. --go_out=. user_service.proto
