.PHONY: gen

gen:
	mkdir -p ./gen/go/feed
	protoc --go_out=./gen/go/feed --go_opt=paths=source_relative \
	--go-grpc_out=./gen/go/feed --go-grpc_opt=paths=source_relative \
	./news_service.proto