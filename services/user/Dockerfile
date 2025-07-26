FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/user/cmd/ ./services/user/cmd/
COPY services/user/internal/ ./services/user/internal/
COPY shared ./shared/
COPY proto ./proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/user ./services/user/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/user .
COPY --from=builder /app/services/user/internal/migrations ./internal/migrations

EXPOSE 50051

CMD ["./user"]
