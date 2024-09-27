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
