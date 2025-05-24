FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV TZ=Europe/Moscow

RUN apk add --no-cache tzdata
RUN ln -sf /usr/share/zoneinfo/Europe/Moscow /etc/localtime

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o bot_question ./cmd/main.go
RUN mkdir -p /etc && cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime

FROM alpine:latest

WORKDIR /root/

ENV TZ=Europe/Moscow


COPY --from=builder /app/bot_question .
COPY --from=builder /etc/localtime /etc/localtime

CMD ["./bot_question"]