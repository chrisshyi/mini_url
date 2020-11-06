run:
	swag init -g ./cmd/web/main.go
	go run ./cmd/web

build:
	swag init -g ./cmd/web/main.go
	go build -o app ./cmd/web