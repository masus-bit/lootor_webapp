FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

RUN apk add --no-cache ca-certificates tzdata

ENV PORT=5111
ENV GIN_MODE=release

EXPOSE 5111

CMD ["/app/main"]