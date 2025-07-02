FROM golang:1.24

WORKDIR /app

COPY . .

WORKDIR /app/services/group

RUN go build -o group ./cmd/main.go

CMD ["./group"]
