.PHONY: gen

gen:
	mkdir -p ./gen/go/microservices
	protoc --go_out=./gen/go/microservices --go_opt=paths=source_relative \
	--go-grpc_out=./gen/go/microservices --go-grpc_opt=paths=source_relative \
	./news_service.proto ./likes.proto ./comments.proto ./notifications.proto