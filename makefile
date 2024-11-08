gen:
	protoc --proto_path=proto --go_out=proto --go-grpc_out=proto proto/*.proto
cli:
	go run client/client.go -p 8002   	
serv:
	go run server/server.go -p 8002   