FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/qr/cmd/ ./services/qr/cmd/
COPY services/qr/internal/ ./services/qr/internal/
COPY shared ./shared/
COPY proto ./proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/qr ./services/qr/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/qr .

EXPOSE 50054

CMD ["./qr"]
