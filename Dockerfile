FROM golang:1.24-alpine
WORKDIR /aplicatie
COPY go.mod ./
COPY . .
RUN go build -o main .
CMD ["./main"]