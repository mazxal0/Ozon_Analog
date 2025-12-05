FROM golang:1.24-alpine AS builder

WORKDIR /app

# Устанавливаем git (иногда нужен для go get)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/app

# ---- Финальный образ ----
FROM alpine:latest

WORKDIR /app

# Важно: сертификаты для HTTPS, MinIO, Google SMTP
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD ["./main"]
