generate-proto-metadata:
	protoc --go_out=./metadata-example --go-grpc_out=./metadata-example ./metadata-example/metadata.proto

generate-proto-interceptor:
	protoc --go_out=./interceptor-example --go-grpc_out=./interceptor-example ./interceptor-example/interceptor.proto

run-server-metadata:
	go run ./metadata-example/server/.

run-client-metadata:
	go run ./metadata-example/client/.

run-server-interceptor:
	go run ./interceptor-example/server/.

run-client-interceptor:
	go run ./interceptor-example/client/.