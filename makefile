run: 
	go run cmd/api/main.go

lint: 
	golangci-lint run

generate:
	go generate ./...

