SRC_FILES=$(wildcard cmd/*)
TEST_BIN=test

create:
	for dir in ${SRC_FILES}; do go build -o $(subst cmd,build,$$dir) $$dir/main.go; done

test:
	go run cmd/${TEST_BIN}/main.go
