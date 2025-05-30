FROM golang:1.23.4 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

RUN apk add --no-cache sqlite

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""

VOLUME /data

EXPOSE 7540

CMD ["./scheduler"]