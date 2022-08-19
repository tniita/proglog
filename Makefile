compile:
	protoc api/v1/*.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=.

test:
	go test -race ./...