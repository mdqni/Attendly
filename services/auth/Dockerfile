FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/auth/cmd/ ./services/auth/cmd/
COPY services/auth/internal/ ./services/auth/internal/
COPY shared ./shared/
COPY proto ./proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/auth ./services/auth/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/auth .
COPY --from=builder /app/services/auth/internal/migrations ./internal/migrations

EXPOSE 50050

CMD ["./auth"]
