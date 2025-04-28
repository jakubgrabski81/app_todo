# Build stage
FROM golang:1.24-alpine AS builder

# Zainstaluj zależności cgo
RUN apk add --no-cache \
    gcc \
    g++ \
    make \
    sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o todo-app

# Run stage
FROM alpine:latest

# Zainstaluj sqlite (bo może być potrzebny w czasie działania aplikacji)
RUN apk add --no-cache sqlite

WORKDIR /root/

COPY --from=builder /app/todo-app .
COPY --from=builder /app/public ./public

# Ustawiamy port
EXPOSE 3001

# Uruchamiamy aplikację
CMD ["./todo-app"]
