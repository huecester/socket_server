SRC_FILES=$(wildcard cmd/*)

create:
	for dir in ${SRC_FILES}; do go build -o $(subst cmd,build,$$dir) $$dir/main.go; done

run:
	go run cmd/server/main.go
