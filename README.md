go mod download

swag init -g cmd/api/main.go --output docs

go run cmd/api/main.go
