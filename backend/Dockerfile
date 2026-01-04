# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/main .

# Создаём директорию для результатов
RUN mkdir -p /app/results

# Устанавливаем timezone
ENV TZ=Europe/Moscow

EXPOSE 8080

CMD ["./main"]
