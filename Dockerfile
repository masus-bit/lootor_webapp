FROM golang:1.24-alpine

RUN apk add --no-cache \
    gcc \
    musl-dev \
    libwebp-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY news_service.proto ./

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN mkdir -p ./gen/go/feed

RUN protoc \
    --go_out=./gen/go/feed \
    --go_opt=paths=source_relative \
    --go-grpc_out=./gen/go/feed \
    --go-grpc_opt=paths=source_relative \
    ./news_service.proto

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/

ENV PORT=5111
EXPOSE 5111
EXPOSE 50051
CMD ["./main"]