FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/group ./services/group
COPY internal/shared ./internal/shared

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags netgo -o bin/group ./services/group/cmd/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/bin/group .
COPY --from=builder /app/services/group/internal/migrations ./internal/migrations

EXPOSE 50052

CMD ["./group"]
