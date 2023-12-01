
run: 
	go run cmd/api/main.go

lint: 
	golangci-lint run

generate:
	go generate ./...

post:
	curl -X POST -H "Content-Type: application/json" -d '{"description": "foo","amount": "123.4567"}' http://localhost:8000/purchase 

post_pretty:
	curl -X POST -H "Content-Type: application/json" -d '{"description": "foo","amount": "123.4567"}' http://localhost:8000/purchase | jq .
