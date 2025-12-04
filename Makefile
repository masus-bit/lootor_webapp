.PHONY: gen

gen:
	mkdir -p ./gen/go/microservices
	protoc --go_out=./gen/go/microservices --go_opt=paths=source_relative \
	--go-grpc_out=./gen/go/microservices --go-grpc_opt=paths=source_relative \
	./protobuf/news_service.proto ./protobuf/likes.proto ./protobuf/comments.proto ./protobuf/notifications.proto ./protobuf/posts.proto ./protobuf/events.proto ./protobuf/tags.proto ./protobuf/photos.proto ./protobuf/achievements.proto