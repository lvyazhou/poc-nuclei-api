# poc-osint-api



## Getting started
```
grpc golang准备工作
1 安装protoc

https://github.com/protocolbuffers/protobuf/releases/tag/v3.14.0
2 安装protoc-gen-go

go get -u github.com/golang/protobuf/protoc-gen-go
3 安装grpc

go get google.golang.org/grpc

go get -u github.com/golang/protobuf/protoc-gen-go@v1.5.2


4 执行proto文件
生成pb和grpc两个文件
protoc --go_out=./ test.proto
protoc --go-grpc_out=./ test.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative test.proto

https://blog.csdn.net/ethunsex/article/details/126697245

https://grpc.io/docs/languages/go/basics/

```