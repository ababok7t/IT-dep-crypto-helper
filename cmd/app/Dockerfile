FROM golang:1.23 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.* .
COPY . .
RUN go mod download

# Копируем исходный код
COPY cmd/ cmd/
COPY /internal/ internal/

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["go","run","cmd/app/main.go"]
