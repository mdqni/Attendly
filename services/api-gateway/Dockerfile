FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/api-gateway ./services/api-gateway
COPY proto ./proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags netgo -ldflags="-s -w" \
    -o bin/api-gateway ./services/api-gateway/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]
