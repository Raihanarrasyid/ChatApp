debug:
	air

run:
	go run cmd/server/main.go

tidy:
	go mod tidy

build-app:
	go build -o ./build/main cmd/main.go

docs:
	swag fmt ./
	swag init -o ./docs -g init.go -d ./internal/app,./internal/controller,./internal/http/request,./internal/http/response,./internal/http/server
