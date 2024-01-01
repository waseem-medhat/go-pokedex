fmt:
	go fmt ./...

lint: fmt
	golangci-lint run
