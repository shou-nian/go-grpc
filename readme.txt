需要在settings - protocol buffers中设置protoc路径：D:\tools\protoc-25.2-win64\include

// 生成service
protoc --grpc-gateway_out=. --go-grpc_out=. --go_out=. user_service.proto

// 生成含swagger文档 ps：需要去protoc-gen-swagger文件夹下执行go build和go install，并把生成的执行文件放入GOPATH和GO ROOT
protoc --swagger_out=. --swagger_opt logtostderr=true --grpc-gateway_out=. --go-grpc_out=. --go_out=. user_service.proto
