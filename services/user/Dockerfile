FROM golang:1.24

WORKDIR /app

COPY . .

WORKDIR /app/services/user

RUN go build -o user ./cmd/main.go

CMD ["./user"]
