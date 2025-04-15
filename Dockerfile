FROM golang:1.24-alpine

RUN apk add --no-cache \
    gcc \
    musl-dev \
    libwebp-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/

ENV PORT=5111
EXPOSE 5111
CMD ["./main"]