FROM golang:1.24-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY comun/ ./comun/
COPY bot/ ./bot/
WORKDIR /app/bot
RUN go mod download
RUN go build -o main .
CMD ["./main"]