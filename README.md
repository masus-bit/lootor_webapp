# Lootor

## Локальная установка


```shell
go mod download
```

## Сгенерировать OpenAPI документацию к API

```shell
swag init -g cmd/api/main.go --output docs
```

## Установка компилятора protoс

```shell
brew install protobuf  ## Для macOS
sudo apt install protobuf-compiler ## Для Linux
choco install protoc ## Для Windows
```

## Генерация protobuf для сервиса новостей

```shell
make gen
```

## Запуск

```shell
go run cmd/api/main.go
```

## Contributing
Хуютинг

## Лицензия

[MIT](https://choosealicense.com/licenses/mit/)