FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/group/cmd/ ./services/group/cmd/
COPY services/group/internal/ ./services/group/internal/
COPY shared ./shared/
COPY proto ./proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/group ./services/group/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/group .
COPY --from=builder /app/services/group/internal/migrations ./internal/migrations

EXPOSE 50052

CMD ["./group"]
