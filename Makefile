test:
	go test ./...
build: test
	go build gocommit.go
