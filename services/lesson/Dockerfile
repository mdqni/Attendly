FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/lesson/cmd/ ./services/lesson/cmd/
COPY services/lesson/internal/ ./services/lesson/internal/
COPY shared ./shared/
COPY proto ./proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/lesson ./services/lesson/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/lesson .
COPY --from=builder /app/services/lesson/internal/migrations ./internal/migrations

EXPOSE 50052

CMD ["./lesson"]
