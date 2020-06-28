build:
	GOOS=linux \
	go build -o dist/routerscan ./cmd/cli/
	cp pkg/librouter/liblibrouter.so dist/