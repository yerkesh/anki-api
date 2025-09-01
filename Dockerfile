# ---------- Build stage ----------
FROM golang:1.24.1 AS builder

# Включаем go modules и отключаем CGO для статичной сборки
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Копируем только go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем остальной код
COPY . .

# Собираем бинарь
RUN go build -o server ./cmd/app

# ---------- Run stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server .

# Запускаем бинарь
ENTRYPOINT ["./server"]
