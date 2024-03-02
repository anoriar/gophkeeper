
## Паспорт сервиса:

* cmd/server/main.go - rest-сервер на Golang (API методы)
* cmd/client/main.go - консольный клиент

## Запуск проекта:

1. В Goland Add Configuration -> go build
2. Run kind = Directory; Directory = к значению, что ide прописало автоматически, надо добавить ```/cmd/shortener```
3. ENVIRONMENT скопировать из ```.env.server-example```

# Компиляция проекта с версией, датой и коммитом
В папке cmd/shortener
go build -ldflags="-X main.buildVersion=1.0.0 -X main.buildDate=$(date -u '+%Y-%m-%dT%H:%M:%S') -X main.buildCommit=abc123" .

# Убрать лишние импорты + gofmt
goimports -local github.com/anoriar/gophkeeper -w ./

## Godoc

godoc -http=:8080

Перейти на 
http://localhost:8080/pkg/github.com/anoriar/shortener/?m=all

# Генерация заглушек godoc

Команда godoc-generate

Либа go install github.com/DimitarPetrov/godoc-generate@latest