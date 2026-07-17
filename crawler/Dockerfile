FROM golang:1.26-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY comun/ ./comun/
COPY crawler/ ./crawler/
WORKDIR /app/crawler
RUN go mod download
RUN go build -o main .
CMD ["./main"]