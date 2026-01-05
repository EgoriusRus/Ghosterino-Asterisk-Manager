# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Копируем исходный код
COPY . .
RUN go mod tidy

# Собираем все бинарники
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o generator ./cmd/generator

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем бинарники из builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/seed .
COPY --from=builder /app/generator .

# Создаём директорию для результатов
RUN mkdir -p /app/results

# Устанавливаем timezone
ENV TZ=Europe/Moscow

EXPOSE 8080

CMD ["./main"]
