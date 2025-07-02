FROM golang:1.24

WORKDIR /app

COPY . .

WORKDIR /app/services/qr

RUN go build -o qr ./cmd/main.go

CMD ["./qr"]
