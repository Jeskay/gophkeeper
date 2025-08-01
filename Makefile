all: server client

server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

generate_protos: generate_proto_message generate_proto_grpc

generate_proto_message:
	protoc ./api/protos/gophkeeper.proto --go_out=. --go_opt=paths=source_relative

generate_proto_grpc:
	protoc ./api/protos/gophkeeper.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative