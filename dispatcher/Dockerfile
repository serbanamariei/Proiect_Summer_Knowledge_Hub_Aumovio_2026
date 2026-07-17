FROM golang:1.26-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY comun/ ./comun/
COPY dispatcher/ ./dispatcher/
WORKDIR /app/dispatcher
RUN go mod download
RUN go build -o main .
CMD while ! nc -z baza_de_date 5432; do sleep 1; done; ./main