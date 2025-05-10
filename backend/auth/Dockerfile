# Этап сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки

# Рабочая директория
WORKDIR /app1

# 1. Копируем только файлы зависимостей (для кэширования слоя)
COPY /backend/auth/go.mod /backend/auth/go.sum ./

# 2. Копируем исходный код
COPY /backend/auth .

# Собираем приложение (убедитесь, что путь к main.go верный)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app1/webbook ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем runtime-зависимости
#RUN apk add --no-cache ca-certificates

# Копируем бинарник и env-файл
COPY --from=builder /app1/webbook /app1/
COPY backend/books/.env /app1/

# Рабочая директория
WORKDIR /app1

# Команда запуска
CMD ["./webbook"]