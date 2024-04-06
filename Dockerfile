FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 8080
EXPOSE 587

CMD ["./main"]
