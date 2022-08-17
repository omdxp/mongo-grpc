genpb:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/*.proto

server:
	go run server/*.go

client:
	go run client/*.go

.PHONY: genpb server client