# Codecov badge

[![codecov](https://codecov.io/gh/iHamsin/practicum-shortener-service/graph/badge.svg?token=Z6XZR4AYOB)](https://codecov.io/gh/iHamsin/practicum-shortener-service)


[![codecov](https://codecov.io/gh/iHamsin/practicum-shortener-service/graphs/sunburst.svg?token=Z6XZR4AYOB)](https://codecov.io/gh/iHamsin/practicum-shortener-service)


# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## DEV запуск

SERVER_ADDRESS="localhost:7070" BASE_URL="http://localhost:7070/env-" go run ./cmd/shortener/main.go -a "localhost:8090" -b "http://localhost:8090/abc-" -f "./temp.db"

SERVER_ADDRESS="localhost:7070" BASE_URL="http://localhost:7070/env-" FILE_STORAGE_PATH="./ENVDB.db" go run ./cmd/shortener/main.go -a "localhost:8090" -b "http://localhost:8090/abc-"

SERVER_ADDRESS="localhost:7070" BASE_URL="http://localhost:7070/env-" FILE_STORAGE_PATH="" go run ./cmd/shortener/main.go -a "localhost:8090" -b "http://localhost:8090/abc-"

SERVER_ADDRESS="localhost:7070" BASE_URL="http://localhost:7070/env-" go run ./cmd/shortener/main.go -a "localhost:8090" -b "http://localhost:8090/abc-" -f ""

SERVER_ADDRESS="localhost:7070" BASE_URL="http://localhost:7070/env-" go run ./cmd/shortener/main.go -a "localhost:8090" -b "http://localhost:8090/abc-" -f "" -d "host=localhost user=yp password=passw0rd sslmode=disable"

SERVER_ADDRESS="localhost:7070" BASE_URL="http://localhost:7070/env-" go run ./cmd/shortener/main.go -d "host=localhost user=yp password=passw0rd sslmode=disable"

## Локальный запуск GitHub Actions

https://github.com/Yandex-Practicum/go-autotests
На github используется v0.10.2 shortenertestbeta-darwin-arm64
https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.10.2/shortenertestbeta-darwin-arm64

chmod +x shortenertest 

## Запуск конкретного action

```
act -j build --container-architecture linux/amd64
act -j golangci --container-architecture linux/amd64
```

## запуск golangci-lint
```
./tests/golangci-lint run ./...
```

## запуск statictest

```
go vet -vettool=./tests/statictest  ./...
```

## запуск всех автотестов от YP

```
go build -o ./tests/ ./cmd/shortener &&
./tests/shortenertest -test.v -binary-path=./tests/shortener -source-path=./ -file-storage-path=./tests/main.db -server-port=8080 -database-dsn="host=localhost user=yp password=passw0rd sslmode=disable"
```

## запуск конкретного автотеста от YP

```
go build -o ./tests/ ./cmd/shortener &&
./tests/shortenertest -test.v -test.run=^TestIteration10$ -binary-path=./tests/shortener -source-path=./ -file-storage-path=./tests/main.db -server-port=8080
```

```
go build -o ./tests/ ./cmd/shortener &&
./tests/shortenertest -test.v -test.run=^TestIteration11$ -binary-path=./tests/shortener -source-path=./ -file-storage-path=./tests/main.db -server-port=8080 -database-dsn="host=localhost user=yp password=passw0rd sslmode=disable database=yp"
```