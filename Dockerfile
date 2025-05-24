# Этап сборки (используем официальный образ Go)
FROM golang:1.23.4 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

# Этап запуска (используем минимальный образ alpine)
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и статические файлы
COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

# Устанавливаем зависимости для SQLite
RUN apk add --no-cache sqlite

# Переменные окружения по умолчанию
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""

# Создаем том для базы данных
VOLUME /data

EXPOSE 7540

CMD ["./scheduler"]