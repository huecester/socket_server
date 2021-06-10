BIN_NAME=server

create:
	go build -o build/${BIN_NAME} cmd/server/main.go

run:
	go run cmd/server/main.go
