SRC_FILES=$(wildcard cmd/*)
TEST_BIN=echo_dual

create:
	for dir in ${SRC_FILES}; do go build -o $$(echo $$dir | sed -e 's/cmd/build/g') $$dir/main.go; done

test:
	go run cmd/${TEST_BIN}/main.go
