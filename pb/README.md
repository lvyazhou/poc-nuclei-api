protoc --go_out=./ test.proto
protoc --go-grpc_out=./ test.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative test.proto