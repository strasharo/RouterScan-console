build:
	LD_LIBRARY_PATH=./pkg/librouter \
	GOOS=linux \
	go build -a  -o routerscan ./cmd/cli/