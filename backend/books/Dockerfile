# Этап сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки

# Рабочая директория
WORKDIR /app

# 1. Копируем только файлы зависимостей (для кэширования слоя)
COPY /backend/books/go.mod /backend/books/go.sum ./

# 2. Копируем исходный код
COPY /backend/books .

# Собираем приложение (убедитесь, что путь к main.go верный)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/webbook ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем runtime-зависимости
#RUN apk add --no-cache ca-certificates

# Копируем бинарник и env-файл
COPY --from=builder /app/webbook /app/
COPY backend/books/.env /app/

# Рабочая директория
WORKDIR /app

# Команда запуска
CMD ["./webbook"]